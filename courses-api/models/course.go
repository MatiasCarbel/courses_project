package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Course representa la estructura de un curso en MongoDB
type Course struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Title          string            `bson:"title" json:"title"`
    Description    string            `bson:"description" json:"description"`
    Category       string            `bson:"category" json:"category"`
    Instructor     string            `bson:"instructor" json:"instructor"`
    Duration       int               `bson:"duration" json:"duration"`
    AvailableSeats int              `bson:"available_seats" json:"available_seats"`
    ImageURL       string            `bson:"image_url" json:"image_url"`
}
