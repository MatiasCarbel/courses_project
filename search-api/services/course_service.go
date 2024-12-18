package services

import (
	"search-api/domain"
	"search-api/repositories"
)

type CourseService struct {
	repo *repositories.SolrRepository
}

func NewCourseService(repo *repositories.SolrRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) UpdateCourse(course domain.Course) error {
	return s.repo.UpdateCourse(course)
}

func (s *CourseService) SearchCourses(query, category, available string) (map[string]interface{}, error) {
	return s.repo.SearchCourses(query, category, available)
}

func (s *CourseService) DeleteCourse(courseID string) error {
	return s.repo.DeleteCourse(courseID)
}
