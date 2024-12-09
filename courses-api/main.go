package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"courses-api/handlers"
	"courses-api/middlewares"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set up MongoDB connection with retry mechanism
	var client *mongo.Client
	var err error
	
	for i := 0; i < 5; i++ {
		clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err = mongo.Connect(ctx, clientOptions)
		cancel()
		
		if err == nil {
			// Test the connection
			err = client.Ping(context.Background(), nil)
			if err == nil {
				break
			}
		}
		
		log.Printf("Failed to connect to MongoDB (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB after 5 attempts: %v", err)
	}
	
	// Ensure disconnection when the main function returns
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Crear el router
	r := mux.NewRouter()

	// Add CORS middleware
	r.Use(middlewares.CorsMiddleware)

	// Rutas públicas
	r.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllCourses(client, w, r)
	}).Methods("GET", "OPTIONS")

	r.HandleFunc("/courses/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCourse(client, w, r)
	}).Methods("GET", "OPTIONS")

	r.HandleFunc("/courses/availability", func(w http.ResponseWriter, r *http.Request) {
		handlers.CheckAvailability(client, w, r)
	}).Methods("POST", "OPTIONS")

	// Rutas protegidas por permisos de administrador
	r.HandleFunc("/courses", middlewares.VerifyAdmin(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCourse(client, w, r)
	})).Methods("POST")

	r.HandleFunc("/courses/{id}", middlewares.VerifyAdmin(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateCourse(client, w, r)
	})).Methods("PUT")

	r.HandleFunc("/courses/{id}", middlewares.VerifyAdmin(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCourse(client, w, r)
	})).Methods("DELETE")

	r.HandleFunc("/enrollments", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateEnrollment(client, w, r)
	}).Methods("POST")

	r.HandleFunc("/enrollments/check/{course_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.CheckEnrollment(client, w, r)
	}).Methods("GET")

	r.HandleFunc("/user/courses/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserCourses(client, w, r)
	}).Methods("GET")

	// Iniciar el servidor
	log.Println("Servidor iniciado en el puerto 8002")
	log.Fatal(http.ListenAndServe(":8002", r))
}
