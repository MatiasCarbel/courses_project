package main

import (
	"context"
	"log"
	"net/http"

	"courses-api/handlers"
	"courses-api/middlewares"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Crear el router
	r := mux.NewRouter()

	// Rutas públicas
	r.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllCourses(client, w, r)
	}).Methods("GET")

	r.HandleFunc("/courses/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCourse(client, w, r)
	}).Methods("GET")

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

	r.HandleFunc("/courses/availability", func(w http.ResponseWriter, r *http.Request) {
		handlers.CalculateAvailability(client, w, r)
	}).Methods("POST")

	r.HandleFunc("/enrollments/check/{course_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.CheckEnrollment(client, w, r)
	}).Methods("GET")

	// Iniciar el servidor
	log.Println("Servidor iniciado en el puerto 8002")
	log.Fatal(http.ListenAndServe(":8002", r))
}
