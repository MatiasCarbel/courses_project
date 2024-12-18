package mongodb

import (
	"context"
	"courses-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnrollmentRepository struct {
	db *mongo.Database
}

func NewEnrollmentRepository(db *mongo.Database) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) collection() *mongo.Collection {
	return r.db.Collection("enrollments")
}

func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	result, err := r.collection().InsertOne(ctx, enrollment)
	if err != nil {
		return models.ErrDatabaseOperation
	}
	enrollment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *EnrollmentRepository) FindByUserID(ctx context.Context, userID int) ([]models.Enrollment, error) {
	cursor, err := r.collection().Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, models.ErrDatabaseOperation
	}
	defer cursor.Close(ctx)

	var enrollments []models.Enrollment
	if err = cursor.All(ctx, &enrollments); err != nil {
		return nil, models.ErrDatabaseOperation
	}
	return enrollments, nil
}

func (r *EnrollmentRepository) CheckEnrollment(ctx context.Context, courseID primitive.ObjectID, userID int) (bool, error) {
	count, err := r.collection().CountDocuments(ctx, bson.M{
		"course_id": courseID,
		"user_id":   userID,
	})
	if err != nil {
		return false, models.ErrDatabaseOperation
	}
	return count > 0, nil
} 