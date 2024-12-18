package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseRepository interface {
	Create(ctx context.Context, course *Course) error
	FindAll(ctx context.Context) ([]Course, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Course, error)
	Update(ctx context.Context, course *Course) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	CheckAvailability(ctx context.Context, courseIDs []primitive.ObjectID) (map[string]int, error)
	FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]Course, error)
}

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *Enrollment) error
	FindByUserID(ctx context.Context, userID int) ([]Enrollment, error)
	CheckEnrollment(ctx context.Context, courseID primitive.ObjectID, userID int) (bool, error)
} 