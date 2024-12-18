package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"search-api/services"
)

type SearchController struct {
	service *services.CourseService
}

func NewSearchController(service *services.CourseService) *SearchController {
	return &SearchController{service: service}
}

func (c *SearchController) SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	available := r.URL.Query().Get("available")

	result, err := c.service.SearchCourses(query, available)
	if err != nil {
		log.Printf("Error searching courses: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
