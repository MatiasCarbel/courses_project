package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

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
	updateURL := fmt.Sprintf("%s/solr/courses/update?commit=true", r.SolrURL)

	// Create a Solr document
	solrDoc := map[string]interface{}{
		"add": map[string]interface{}{
			"doc": course,
		},
	}

	courseJSON, err := json.Marshal(solrDoc)
	if err != nil {
		return fmt.Errorf("error marshaling course: %w", err)
	}

	req, err := http.NewRequest("POST", updateURL, bytes.NewBuffer(courseJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error updating Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Solr update failed: %s, response: %s", resp.Status, string(body))
	}

	return nil
}

func (r *SolrRepository) SearchCourses(query, available string) (map[string]interface{}, error) {
	searchURL := fmt.Sprintf("%s/solr/courses/select?wt=json", r.SolrURL)

	if query == "" || query == "*:*" {
		searchURL += "&q=*:*"
	} else {
		searchURL += fmt.Sprintf("&q=title:%s OR description:%s",
			url.QueryEscape(query), url.QueryEscape(query))
	}

	if available == "true" {
		searchURL += "&fq=available_seats:[1 TO *]"
	}

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error querying Solr: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding Solr response: %w", err)
	}

	return result, nil
}
