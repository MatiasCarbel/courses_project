package mongodb

import (
	"context"
	"courses-api/models"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseRepository struct {
	db *mongo.Database
}

func NewCourseRepository(db *mongo.Database) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) collection() *mongo.Collection {
	return r.db.Collection("courses")
}

func (r *CourseRepository) Create(ctx context.Context, course *models.Course) error {
	// Check for duplicates
	filter := bson.M{"title": course.Title, "instructor": course.Instructor}
	count, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("a course with the same title and instructor already exists")
	}

	result, err := r.collection().InsertOne(ctx, course)
	if err != nil {
		return err
	}
	course.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *CourseRepository) FindAll(ctx context.Context) ([]models.Course, error) {
	cursor, err := r.collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var courses []models.Course
	if err = cursor.All(ctx, &courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *CourseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Course, error) {
	var course models.Course
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrCourseNotFound
		}
		return nil, models.ErrDatabaseOperation
	}
	return &course, nil
}

func (r *CourseRepository) Update(ctx context.Context, course *models.Course) error {
	filter := bson.M{"_id": course.ID}
	update := bson.M{"$set": course}
	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("course not found")
	}
	return nil
}

func (r *CourseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("course not found")
	}
	return nil
}

func (r *CourseRepository) CheckAvailability(ctx context.Context, courseIDs []primitive.ObjectID) (map[string]int, error) {
	filter := bson.M{
		"_id": bson.M{"$in": courseIDs},
		"available_seats": bson.M{"$gt": 0},
	}

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	availability := make(map[string]int)
	var course models.Course
	for cursor.Next(ctx) {
		if err := cursor.Decode(&course); err != nil {
			return nil, err
		}
		availability[course.ID.Hex()] = course.AvailableSeats
	}

	return availability, nil
}

func (r *CourseRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Course, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, models.ErrDatabaseOperation
	}
	defer cursor.Close(ctx)

	var courses []models.Course
	if err = cursor.All(ctx, &courses); err != nil {
		return nil, models.ErrDatabaseOperation
	}
	return courses, nil
}