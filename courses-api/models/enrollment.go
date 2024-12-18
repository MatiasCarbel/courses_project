package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Enrollment struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    CourseID  primitive.ObjectID `bson:"course_id" json:"course_id"`
    UserID    int               `bson:"user_id" json:"user_id"`
    Date      time.Time         `bson:"date" json:"date"`
} 