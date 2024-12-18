package main

import (
	"log"
	"net/http"

	"search-api/controllers"
	"search-api/repositories"
	"search-api/services"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

func main() {
	// Initialize Solr Repository
	solrRepo := repositories.NewSolrRepository()

	// Initialize Services
	courseService := services.NewCourseService(solrRepo)
	rabbitMQService := services.NewRabbitMQService(courseService)

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Start consuming RabbitMQ messages
	go rabbitMQService.ConsumeRabbitMQMessages(conn)

	// Initialize Controllers
	searchController := controllers.NewSearchController(courseService)

	// Set up HTTP routes
	r := mux.NewRouter()
	r.HandleFunc("/search", searchController.SearchHandler).Methods("GET")

	log.Println("Search API running on port 8003")
	log.Fatal(http.ListenAndServe(":8003", r))
}
