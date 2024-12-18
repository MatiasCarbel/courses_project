package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"search-api/domain"
)

type SolrRepository struct {
	SolrURL string
}

func NewSolrRepository() *SolrRepository {
	solrURL := os.Getenv("SOLR_URL")
	if solrURL == "" {
		solrURL = "http://solr:8983"
	}
	return &SolrRepository{SolrURL: solrURL}
}

func (r *SolrRepository) UpdateCourse(course domain.Course) error {
	updateURL := fmt.Sprintf("%s/solr/courses/update/json/docs?commit=true", r.SolrURL)
	
	doc := map[string]interface{}{
		"id":              course.ID,
		"title":          course.Title,
		"description":    course.Description,
		"instructor":     course.Instructor,
		"category":       course.Category,
		"image_url":      course.ImageURL,
		"duration":       course.Duration,
		"available_seats": course.AvailableSeats,
	}

	jsonData, err := json.Marshal([]interface{}{doc})
	if err != nil {
		return fmt.Errorf("error marshaling update data: %w", err)
	}

	req, err := http.NewRequest("POST", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Solr update failed: %s, response: %s", resp.Status, string(body))
	}

	return nil
}

func (r *SolrRepository) SearchCourses(query, category, available string) (map[string]interface{}, error) {
	searchURL := fmt.Sprintf("%s/solr/courses/select", r.SolrURL)
	
	// Build query parameters
	params := url.Values{}
	params.Add("wt", "json") // Explicitly request JSON response
	
	// Base query
	baseQuery := "*:*"
	
	// Handle text search if present
	if query != "" && query != "*:*" {
		baseQuery = fmt.Sprintf("title:*%s* OR description:*%s*", query, query)
	}
	
	// Add base query
	params.Add("q", baseQuery)
	
	// Handle category filter
	if category != "" {
		// Use filter query for exact category matching
		params.Add("fq", fmt.Sprintf("category:\"%s\"", category))
	}
	
	// Handle available seats filter
	if available == "true" {
		params.Add("fq", "available_seats:[1 TO *]")
	}
	
	// Specify fields to return
	params.Add("fl", "id,title,description,instructor,category,image_url,duration,available_seats")
	
	// Add parameters to URL
	finalURL := fmt.Sprintf("%s?%s", searchURL, params.Encode())
	
	log.Printf("Solr query URL: %s", finalURL)
	
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying Solr: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Solr query failed: %s, response: %s", resp.Status, string(body))
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding Solr response: %w", err)
	}
	
	return result, nil
}

func (r *SolrRepository) DeleteCourse(courseID string) error {
	deleteURL := fmt.Sprintf("%s/solr/courses/update?commit=true", r.SolrURL)
	deleteQuery := fmt.Sprintf(`{"delete": { "id": "%s" }}`, courseID)

	req, err := http.NewRequest("POST", deleteURL, strings.NewReader(deleteQuery))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete course from Solr: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
