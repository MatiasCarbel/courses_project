package main

import (
	"courses-api/config"
	"courses-api/controllers"
	"courses-api/middlewares"
	"courses-api/repositories/mongodb"
	"courses-api/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize MongoDB connection
	db, err := config.InitMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	courseRepo := mongodb.NewCourseRepository(db)
	enrollmentRepo := mongodb.NewEnrollmentRepository(db)
	
	// Initialize message queue
	messageQueue := services.NewRabbitMQService()
	
	// Initialize services
	courseService := services.NewCourseService(courseRepo, enrollmentRepo, messageQueue)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo, courseRepo, messageQueue)
	
	// Initialize controllers
	courseController := controllers.NewCourseController(courseService)
	enrollmentController := controllers.NewEnrollmentController(enrollmentService)
	serviceController := controllers.NewServiceController()

	// Set up router
	r := mux.NewRouter()
	r.Use(middlewares.CorsMiddleware)

	// Register routes
	r.HandleFunc("/courses", courseController.GetAllCourses).Methods("GET", "OPTIONS")
	r.HandleFunc("/courses/myCourses", middlewares.VerifyToken(courseController.GetUserCourses)).Methods("GET", "OPTIONS")
	r.HandleFunc("/courses/{id}", courseController.GetCourse).Methods("GET", "OPTIONS")
	r.HandleFunc("/courses/availability", courseController.CheckAvailability).Methods("POST", "OPTIONS")
	
	// Protected routes
	r.HandleFunc("/courses", middlewares.VerifyAdmin(courseController.CreateCourse)).Methods("POST")
	r.HandleFunc("/courses/{id}", middlewares.VerifyAdmin(courseController.UpdateCourse)).Methods("PUT")
	r.HandleFunc("/courses/{id}", middlewares.VerifyAdmin(courseController.DeleteCourse)).Methods("DELETE")

	// Enrollment routes
	r.HandleFunc("/enrollments", middlewares.VerifyToken(enrollmentController.CreateEnrollment)).Methods("POST")
	r.HandleFunc("/enrollments/check/{courseId}", middlewares.VerifyToken(enrollmentController.CheckEnrollment)).Methods("GET")

	// Service routes
	r.HandleFunc("/api/services", middlewares.VerifyAdmin(serviceController.GetServices)).Methods("GET")
	r.HandleFunc("/api/services", middlewares.VerifyAdmin(serviceController.AddInstance)).Methods("POST")
	r.HandleFunc("/api/services/{id}", middlewares.VerifyAdmin(serviceController.RemoveInstance)).Methods("DELETE")

	log.Println("Server started on port 8002")
	log.Fatal(http.ListenAndServe(":8002", r))
}
