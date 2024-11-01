package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type Course struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Instructor     string `json:"instructor"`
	Duration       int    `json:"duration"`
	AvailableSeats int    `json:"available_seats"`
}

func main() {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ { // Retry 5 times
			conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
			if err == nil {
					break
			}
			log.Printf("Failed to connect to RabbitMQ: %v. Retrying...", err)
			time.Sleep(5 * time.Second) // Wait before retrying
	}

	if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ after retries: %v", err)
	}
	defer conn.Close()

	// Consume messages from RabbitMQ
	go consumeRabbitMQMessages(conn)

	// Set up HTTP server
	r := mux.NewRouter()
	r.HandleFunc("/search", searchHandler).Methods("GET")

	log.Println("Search API running on port 8003")
	log.Fatal(http.ListenAndServe(":8003", r))
}

func consumeRabbitMQMessages(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
			log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
			"course_updates", // name
			true,            // durable
			false,           // delete when unused
			false,           // exclusive
			false,           // no-wait
			nil,            // arguments
	)
	if err != nil {
			log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
			"course_updates", // queue
			"",              // consumer
			true,            // auto-ack
			false,           // exclusive
			false,           // no-local
			false,           // no-wait
			nil,             // args
	)
	if err != nil {
			log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
			var update map[string]interface{}
			if err := json.Unmarshal(msg.Body, &update); err != nil {
					log.Printf("Error decoding JSON: %v", err)
					continue
			}

			// Check if this is a seat update
			if seatChange, ok := update["seat_change"].(float64); ok {
					handleSeatUpdate(update["course_id"].(string), int(seatChange))
					continue
			}

			// Handle course update/creation
			course := Course{
					ID:             update["id"].(string),
					Title:          update["title"].(string),
					Description:    update["description"].(string),
					Instructor:     update["instructor"].(string),
					Duration:       int(update["duration"].(float64)),
					AvailableSeats: int(update["available_seats"].(float64)),
			}

			updateSolR(course)
	}
}

func handleSeatUpdate(courseID string, seatChange int) {
	solrURL := fmt.Sprintf("http://solr:8983/solr/courses/select?q=id:%s&wt=json", url.QueryEscape(courseID))
	resp, err := http.Get(solrURL)
	if err != nil {
			log.Printf("Error querying SolR: %v", err)
			return
	}
	defer resp.Body.Close()

	var solrResponse struct {
			Response struct {
					Docs []Course `json:"docs"`
			} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&solrResponse); err != nil {
			log.Printf("Error decoding SolR response: %v", err)
			return
	}

	if len(solrResponse.Response.Docs) == 0 {
			log.Printf("Course not found in SolR: %s", courseID)
			return
	}

	// Update available seats
	course := solrResponse.Response.Docs[0]
	course.AvailableSeats += seatChange

	// Update the document in Solr
	updateSolR(course)
}

func updateSolR(course Course) {
	solrURL := "http://solr:8983/solr/courses/update/json/docs?commit=true"
	
	// Create a map for the Solr document
	solrDoc := map[string]interface{}{
			"id":              course.ID,
			"title":          course.Title,
			"description":    course.Description,
			"instructor":     course.Instructor,
			"duration":       course.Duration,
			"available_seats": course.AvailableSeats,
	}
	
	courseJSON, err := json.Marshal(solrDoc)
	if err != nil {
			log.Printf("Error marshaling course: %v", err)
			return
	}

	resp, err := http.Post(solrURL, "application/json", bytes.NewBuffer(courseJSON))
	if err != nil {
			log.Printf("Error updating SolR: %v", err)
			return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("SolR update failed with status: %v, body: %s", resp.Status, string(body))
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("q")
	availableOnly := r.URL.Query().Get("available") == "true"

	// Construct SolR query URL
	solrURL := fmt.Sprintf("http://solr:8983/solr/courses/select?q=%s&wt=json", url.QueryEscape(query))

	// Make request to SolR
	resp, err := http.Get(solrURL)
	if err != nil {
			log.Printf("Error querying SolR: %v", err)
			http.Error(w, "Error querying SolR", http.StatusInternalServerError)
			return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
			log.Printf("Error reading SolR response: %v", err)
			http.Error(w, "Error reading SolR response", http.StatusInternalServerError)
			return
	}

	if resp.StatusCode != http.StatusOK {
			log.Printf("SolR query failed with status: %v", resp.Status)
			http.Error(w, "SolR query failed", http.StatusInternalServerError)
			return
	}

	// Decode SolR response
	var solrResponse struct {
			Response struct {
					Docs []Course `json:"docs"`
			} `json:"response"`
	}
	if err := json.Unmarshal(body, &solrResponse); err != nil {
			log.Printf("Error decoding SolR response: %v", err)
			http.Error(w, "Error decoding SolR response", http.StatusInternalServerError)
			return
	}

	// Filter courses based on availability
	var filteredCourses []Course
	for _, course := range solrResponse.Response.Docs {
			if !availableOnly || course.AvailableSeats > 0 {
					filteredCourses = append(filteredCourses, course)
			}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredCourses)
}
