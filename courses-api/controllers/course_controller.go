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

type CourseService interface {
	CreateCourse(ctx context.Context, course *models.Course) error
	GetAllCourses(ctx context.Context) ([]models.Course, error)
	GetCourse(ctx context.Context, id primitive.ObjectID) (*models.Course, error)
	UpdateCourse(ctx context.Context, course *models.Course) error
	DeleteCourse(ctx context.Context, id primitive.ObjectID) error
	CheckAvailability(ctx context.Context, courseIDs []primitive.ObjectID) (map[string]int, error)
	GetUserCourses(ctx context.Context, userID int) ([]models.Course, error)
}

type CourseController struct {
	service CourseService
}

func NewCourseController(service CourseService) *CourseController {
	return &CourseController{service: service}
}

func (c *CourseController) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid request body",
		})
		return
	}

	if err := c.service.CreateCourse(r.Context(), &course); err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusCreated,
		Data:   course,
	})
}

func (c *CourseController) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := c.service.GetAllCourses(r.Context())
	if err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusInternalServerError,
			Error:  "Error fetching courses",
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Data:   courses,
	})
}

func (c *CourseController) GetCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid course ID",
		})
		return
	}

	course, err := c.service.GetCourse(r.Context(), courseID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "course not found" {
			status = http.StatusNotFound
		}
		views.JSON(w, views.Response{
			Status: status,
			Error:  err.Error(),
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Data:   course,
	})
}

func (c *CourseController) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid course ID",
		})
		return
	}

	var course models.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid request body",
		})
		return
	}

	course.ID = courseID
	if err := c.service.UpdateCourse(r.Context(), &course); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "course not found" {
			status = http.StatusNotFound
		}
		views.JSON(w, views.Response{
			Status: status,
			Error:  err.Error(),
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Data:   course,
	})
}

func (c *CourseController) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid course ID",
		})
		return
	}

	if err := c.service.DeleteCourse(r.Context(), courseID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "course not found" {
			status = http.StatusNotFound
		}
		views.JSON(w, views.Response{
			Status: status,
			Error:  err.Error(),
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Message: "Course successfully deleted",
	})
}

func (c *CourseController) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var courseIDStrings []string
	if err := json.NewDecoder(r.Body).Decode(&courseIDStrings); err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid request body",
		})
		return
	}

	courseIDs := make([]primitive.ObjectID, 0, len(courseIDStrings))
	for _, idStr := range courseIDStrings {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			views.JSON(w, views.Response{
				Status: http.StatusBadRequest,
				Error:  "Invalid course ID format",
			})
			return
		}
		courseIDs = append(courseIDs, id)
	}

	availability, err := c.service.CheckAvailability(r.Context(), courseIDs)
	if err != nil {
		views.JSON(w, views.Response{
			Status: http.StatusInternalServerError,
			Error:  "Error checking availability",
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Data:   availability,
	})
}

func (c *CourseController) GetUserCourses(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT claims (assuming it's set in the context by middleware)
	userID := r.Context().Value("userID").(int)
	if userID == 0 {
		views.JSON(w, views.Response{
			Status: http.StatusUnauthorized,
			Error:  models.ErrUnauthorized.Error(),
		})
		return
	}

	courses, err := c.service.GetUserCourses(r.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case models.ErrUnauthorized:
			status = http.StatusUnauthorized
		case models.ErrDatabaseOperation:
			status = http.StatusInternalServerError
		default:
			status = http.StatusInternalServerError
		}
		
		views.JSON(w, views.Response{
			Status: status,
			Error:  err.Error(),
		})
		return
	}

	views.JSON(w, views.Response{
		Status: http.StatusOK,
		Data:   courses,
	})
}