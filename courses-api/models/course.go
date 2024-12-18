package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Title          string            `bson:"title" json:"title"`
    Description    string            `bson:"description" json:"description"`
    Instructor     string            `bson:"instructor" json:"instructor"`
    Duration       int               `bson:"duration" json:"duration"`
    AvailableSeats int              `bson:"available_seats" json:"available_seats"`
    Category       string            `bson:"category" json:"category"`
    ImageURL       string            `bson:"image_url" json:"image_url"`
}

var ValidCategories = []string{"web-development", "mobile-development", "data-science", "design", "business"}

func (c *Course) Validate() error {
    if c.Title == "" || c.Description == "" || c.Instructor == "" || 
       c.Duration <= 0 || c.AvailableSeats <= 0 || c.Category == "" {
        return errors.New("all fields are required and numeric values must be greater than 0")
    }

    isValidCategory := false
    for _, validCat := range ValidCategories {
        if c.Category == validCat {
            isValidCategory = true
            break
        }
    }

    if !isValidCategory {
        return errors.New("invalid category")
    }

    return nil
}
