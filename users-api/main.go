package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"users-api/handlers"
	"users-api/repositories"
	"users-api/services"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	log.Println("Conexión exitosa a la base de datos")

	// Initialize repository and service
	userRepo := repositories.NewSQLUserRepository(db)
	userService := services.NewUserService(userRepo, []byte(os.Getenv("JWT_SECRET")))

	// Initialize handlers
	userHandlers := handlers.NewUserHandlers(userService)

	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/users", userHandlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandlers.GetUser).Methods("GET")
	r.HandleFunc("/user/login", userHandlers.LoginUser).Methods("POST")

	log.Println("Iniciando servidor en el puerto 8001")
	log.Fatal(http.ListenAndServe(":8001", r))
}
