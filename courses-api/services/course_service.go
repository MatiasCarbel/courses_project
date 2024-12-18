package services

import (
	"context"
	"courses-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseService struct {
    repo           models.CourseRepository
    enrollmentRepo models.EnrollmentRepository
    messageQueue   MessageQueue
}

func NewCourseService(repo models.CourseRepository, enrollmentRepo models.EnrollmentRepository, mq MessageQueue) *CourseService {
    return &CourseService{
        repo:           repo,
        enrollmentRepo: enrollmentRepo,
        messageQueue:   mq,
    }
}

func (s *CourseService) CreateCourse(ctx context.Context, course *models.Course) error {
    if err := course.Validate(); err != nil {
        return err
    }

    if err := s.repo.Create(ctx, course); err != nil {
        return err
    }

    // Publish course creation event
    return s.messageQueue.PublishCourseUpdate(course, "upsert")
}

func (s *CourseService) GetAllCourses(ctx context.Context) ([]models.Course, error) {
    return s.repo.FindAll(ctx)
}

func (s *CourseService) GetCourse(ctx context.Context, id primitive.ObjectID) (*models.Course, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *CourseService) UpdateCourse(ctx context.Context, course *models.Course) error {
    if err := course.Validate(); err != nil {
        return err
    }

    if err := s.repo.Update(ctx, course); err != nil {
        return err
    }

    // Publish course update event
    return s.messageQueue.PublishCourseUpdate(course, "upsert")
}

func (s *CourseService) DeleteCourse(ctx context.Context, id primitive.ObjectID) error {
    if err := s.repo.Delete(ctx, id); err != nil {
        return err
    }

    // Publish course deletion event
    return s.messageQueue.PublishCourseDelete(id)
}

func (s *CourseService) CheckAvailability(ctx context.Context, courseIDs []primitive.ObjectID) (map[string]int, error) {
    return s.repo.CheckAvailability(ctx, courseIDs)
}

func (s *CourseService) GetUserCourses(ctx context.Context, userID int) ([]models.Course, error) {
    enrollments, err := s.enrollmentRepo.FindByUserID(ctx, userID)
    if err != nil {
        return nil, models.ErrDatabaseOperation
    }

    if len(enrollments) == 0 {
        return []models.Course{}, nil
    }

    var courseIDs []primitive.ObjectID
    for _, enrollment := range enrollments {
        courseIDs = append(courseIDs, enrollment.CourseID)
    }

    return s.repo.FindByIDs(ctx, courseIDs)
} 