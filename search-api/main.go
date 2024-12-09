package main

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
	Category       string `json:"category"`
	ImageURL       string `json:"image_url"`
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
			log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
			"course_updates", // name
			true,            // durable
			false,           // delete when unused
			false,           // exclusive
			false,           // no-wait
			nil,            // arguments
	)
	if err != nil {
			log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
	)
	if err != nil {
			log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
			for d := range msgs {
					log.Printf("Received message from RabbitMQ: %s", string(d.Body))
					
					var event map[string]interface{}
					if err := json.Unmarshal(d.Body, &event); err != nil {
							log.Printf("Error parsing message: %v", err)
							continue
					}

					action := event["action"].(string)
					courseData := event["course"].(map[string]interface{})

					log.Printf("Processing %s action for course ID: %s", action, courseData["id"])

					switch action {
					case "upsert":
							course := Course{
									ID:             courseData["id"].(string),
									Title:          courseData["title"].(string),
									Description:    courseData["description"].(string),
									Instructor:     courseData["instructor"].(string),
									Duration:       int(courseData["duration"].(float64)),
									AvailableSeats: int(courseData["available_seats"].(float64)),
									Category:       courseData["category"].(string),
									ImageURL:       courseData["image_url"].(string),
							}
							updateSolR(course)
							log.Printf("Course %s updated in Solr", course.ID)
					case "delete":
							courseID := courseData["id"].(string)
							deleteSolRDocument(courseID)
							log.Printf("Course %s deleted from Solr", courseID)
					}
			}
	}()

	<-forever
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
	solrURL := "http://solr:8983/solr/courses/update?commit=true"
	
	solrDoc := map[string]interface{}{
		"add": map[string]interface{}{
			"doc": map[string]interface{}{
				"id":              course.ID,
				"title":          course.Title,
				"description":    course.Description,
				"category":       course.Category,
				"instructor":     course.Instructor,
				"duration":       course.Duration,
				"available_seats": course.AvailableSeats,
				"image_url":     course.ImageURL,
			},
		},
	}
	
	courseJSON, err := json.Marshal(solrDoc)
	if err != nil {
		log.Printf("Error marshaling course: %v", err)
		return
	}

	req, err := http.NewRequest("POST", solrURL, bytes.NewBuffer(courseJSON))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error updating SolR: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("SolR update failed with status: %v, body: %s", resp.Status, string(body))
	} else {
		log.Printf("Successfully updated course %s in Solr", course.ID)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")

	solrURL := os.Getenv("SOLR_URL")
	if solrURL == "" {
		solrURL = "http://solr:8983"
	}

	// Build Solr query URL
	searchURL := fmt.Sprintf("%s/solr/courses/select?wt=json", solrURL)
	
	// Build query string
	queryParts := []string{}
	
	if query != "" {
		// Escapar espacios y caracteres especiales
		escapedQuery := strings.Replace(query, " ", "\\ ", -1)
		queryParts = append(queryParts, fmt.Sprintf("(title:*%s* OR description:*%s*)", 
			url.QueryEscape(escapedQuery), url.QueryEscape(escapedQuery)))
	}
	
	if category != "" {
		queryParts = append(queryParts, fmt.Sprintf("category:%s", 
			url.QueryEscape(category)))
	}
	
	finalQuery := "*:*"
	if len(queryParts) > 0 {
		finalQuery = strings.Join(queryParts, " AND ")
	}
	
	searchURL += "&q=" + url.QueryEscape(finalQuery)

	// Execute search request
	resp, err := http.Get(searchURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response to client
	io.Copy(w, resp.Body)
}

func deleteSolRDocument(courseID string) {
	solrURL := "http://solr:8983/solr/courses/update?commit=true"
	
	deleteDoc := map[string]interface{}{
		"delete": courseID,
	}
	
	deleteJSON, err := json.Marshal(deleteDoc)
	if err != nil {
		log.Printf("Error marshaling delete request: %v", err)
		return
	}

	req, err := http.NewRequest("POST", solrURL, bytes.NewBuffer(deleteJSON))
	if err != nil {
		log.Printf("Error creating delete request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error deleting from SolR: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("SolR delete failed with status: %v, body: %s", resp.Status, string(body))
	}
}
