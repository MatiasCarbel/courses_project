package services

import (
	"context"
	"courses-api/models"
)

type EnrollmentService struct {
	enrollmentRepo models.EnrollmentRepository
	courseRepo     models.CourseRepository
	messageQueue   MessageQueue
}

func NewEnrollmentService(
	enrollmentRepo models.EnrollmentRepository,
	courseRepo models.CourseRepository,
	messageQueue MessageQueue,
) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		messageQueue:   messageQueue,
	}
}

func (s *EnrollmentService) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	course, err := s.courseRepo.FindByID(ctx, enrollment.CourseID)
	if err != nil {
		if err == models.ErrCourseNotFound {
			return models.ErrCourseNotFound
		}
		return err
	}

	if course.AvailableSeats <= 0 {
		return models.ErrNoAvailableSeats
	}

	enrolled, err := s.enrollmentRepo.CheckEnrollment(ctx, enrollment.CourseID, enrollment.UserID)
	if err != nil {
		return err
	}
	if enrolled {
		return models.ErrAlreadyEnrolled
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return err
	}

	course.AvailableSeats--
	return s.courseRepo.Update(ctx, course)
}

func (s *EnrollmentService) GetUserEnrollments(ctx context.Context, userID int) ([]models.Enrollment, error) {
	return s.enrollmentRepo.FindByUserID(ctx, userID)
} 