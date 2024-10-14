package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"courses-api/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
    // Datos de conexión obtenidos de variables de entorno
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error abriendo la base de datos: %v", err)
    }
    defer db.Close()

    // Verificamos la conexión
    err = db.Ping()
    if err != nil {
        log.Fatalf("Error conectándose a la base de datos: %v", err)
    }

    log.Println("Conexión exitosa a PostgreSQL")

    // Iniciamos el router para los endpoints
    r := mux.NewRouter()

    // Definir rutas para el microservicio
    r.HandleFunc("/courses", handlers.GetCourses).Methods("GET")
    r.HandleFunc("/courses/{id}", handlers.GetCourseByID).Methods("GET")
    r.HandleFunc("/courses", handlers.CreateCourse).Methods("POST")
    r.HandleFunc("/courses/{id}", handlers.UpdateCourse).Methods("PUT")
    r.HandleFunc("/courses/{id}", handlers.DeleteCourse).Methods("DELETE")

    // Puerto en el que correrá el servicio
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    log.Printf("Iniciando server en el puerto %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
