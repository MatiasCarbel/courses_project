package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"users-api/handlers"
	dao "users-api/repositories"
	"users-api/services"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	log.Println("Conexión exitosa a la base de datos")

	// Initialize the layers
	userDAO := &dao.UserDAO{DB: db}
	userService := &services.UserService{UserDAO: userDAO}
	userHandler := &handlers.UserHandler{UserService: userService}

	r := mux.NewRouter()

	// Definir rutas

	r.HandleFunc("/users", userHandler.CreateUserHandler).Methods("POST")

	r.HandleFunc("/users/{id}", userHandler.GetUserHandler).Methods("GET")

	r.HandleFunc("/users/{id}", userHandler.DeleteUserHandler).Methods("DELETE")

	r.HandleFunc("/users/{id}", userHandler.UpdateUserHandler).Methods("PUT")

	r.HandleFunc("/login", userHandler.LoginUser).Methods("POST")

	log.Println("Iniciando servidor en el puerto 8001")
	log.Fatal(http.ListenAndServe(":8001", r))
}
