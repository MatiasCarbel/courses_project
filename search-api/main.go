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
					var event map[string]interface{}
					if err := json.Unmarshal(d.Body, &event); err != nil {
							log.Printf("Error parsing message: %v", err)
							continue
					}

					action := event["action"].(string)
					courseData := event["course"].(map[string]interface{})

					switch action {
					case "upsert":
							course := Course{
									ID:             courseData["id"].(string),
									Title:          courseData["title"].(string),
									Description:    courseData["description"].(string),
									Instructor:     courseData["instructor"].(string),
									Duration:       int(courseData["duration"].(float64)),
									AvailableSeats: int(courseData["available_seats"].(float64)),
							}
							updateSolR(course)
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
	
	// Create a map for the Solr document
	solrDoc := map[string]interface{}{
		"add": map[string]interface{}{
			"doc": map[string]interface{}{
				"id":              course.ID,
				"title":          course.Title,
				"description":    course.Description,
				"instructor":     course.Instructor,
				"duration":       course.Duration,
				"available_seats": course.AvailableSeats,
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
	available := r.URL.Query().Get("available")

	solrURL := os.Getenv("SOLR_URL")
	if solrURL == "" {
		solrURL = "http://solr:8983"
	}

	// Build Solr query URL
	searchURL := fmt.Sprintf("%s/solr/courses/select?wt=json", solrURL)
	
	if query == "" || query == "*:*" {
		searchURL += "&q=*:*"
	} else {
		searchURL += fmt.Sprintf("&q=title:%s OR description:%s", 
			url.QueryEscape(query), url.QueryEscape(query))
	}

	if available == "true" {
		searchURL += "&fq=available_seats:[1 TO *]"
	}

	log.Printf("Querying Solr: %s", searchURL)

	resp, err := http.Get(searchURL)
	if err != nil {
		log.Printf("Error querying Solr: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding Solr response: %v", err)
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
