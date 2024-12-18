package models

import "errors"

// Course-related errors
var (
    ErrCourseNotFound     = errors.New("course not found")
    ErrDuplicateCourse    = errors.New("a course with the same title and instructor already exists")
    ErrInvalidCourseData  = errors.New("invalid course data")
    ErrInvalidCategory    = errors.New("invalid category")
)

// Enrollment-related errors
var (
    ErrNoAvailableSeats  = errors.New("no available seats in the course")
    ErrAlreadyEnrolled   = errors.New("user is already enrolled in this course")
    ErrEnrollmentNotFound = errors.New("enrollment not found")
)

// Authentication/Authorization errors
var (
    ErrUnauthorized      = errors.New("unauthorized access")
    ErrForbidden         = errors.New("forbidden access")
    ErrInvalidToken      = errors.New("invalid or expired token")
)

// Database errors
var (
    ErrDatabaseConnection = errors.New("database connection error")
    ErrDatabaseOperation  = errors.New("database operation failed")
)

// Message Queue errors
var (
    ErrMessageQueueConnection = errors.New("message queue connection error")
    ErrMessagePublishing     = errors.New("error publishing message")
) 