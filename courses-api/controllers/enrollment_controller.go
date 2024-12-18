package controllers

import (
	"context"
	"courses-api/models"
	"courses-api/views"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnrollmentService interface {
    CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
    GetUserEnrollments(ctx context.Context, userID int) ([]models.Enrollment, error)
    CheckEnrollment(ctx context.Context, courseID primitive.ObjectID, userID int) (bool, error)
}

type EnrollmentController struct {
    service EnrollmentService
}

func NewEnrollmentController(service EnrollmentService) *EnrollmentController {
    return &EnrollmentController{service: service}
}

func (c *EnrollmentController) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
    // Get userID from context (set by auth middleware)
    userID := r.Context().Value("userID").(int)

    var enrollment models.Enrollment
    if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
        views.JSON(w, views.Response{
            Status: http.StatusBadRequest,
            Error:  "Invalid request body",
        })
        return
    }

    // Set the UserID from the JWT claims
    enrollment.UserID = userID

    if err := c.service.CreateEnrollment(r.Context(), &enrollment); err != nil {
        status := http.StatusInternalServerError
        switch err {
        case models.ErrCourseNotFound:
            status = http.StatusNotFound
        case models.ErrNoAvailableSeats:
            status = http.StatusConflict
        case models.ErrAlreadyEnrolled:
            status = http.StatusConflict
        }
        views.JSON(w, views.Response{
            Status: status,
            Error:  err.Error(),
        })
        return
    }

    views.JSON(w, views.Response{
        Status: http.StatusCreated,
        Data:   enrollment,
    })
}

func (c *EnrollmentController) CheckEnrollment(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)
    vars := mux.Vars(r)
    courseID, err := primitive.ObjectIDFromHex(vars["courseId"])
    if err != nil {
        views.JSON(w, views.Response{
            Status: http.StatusBadRequest,
            Error:  "Invalid course ID",
        })
        return
    }

    enrolled, err := c.service.CheckEnrollment(r.Context(), courseID, userID)
    if err != nil {
        views.JSON(w, views.Response{
            Status: http.StatusInternalServerError,
            Error:  "Error checking enrollment",
        })
        return
    }

    views.JSON(w, views.Response{
        Status: http.StatusOK,
        Data: map[string]bool{
            "enrolled": enrolled,
        },
    })
} 