package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtSecret = []byte("uccdemy") // Use a secure key in production

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

// Course representa la estructura de un curso
type Course struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title          string            `bson:"title" json:"title"`
	Description    string            `bson:"description" json:"description"`
	Instructor     string            `bson:"instructor" json:"instructor"`
	Duration       int               `bson:"duration" json:"duration"`
	AvailableSeats int              `bson:"available_seats" json:"available_seats"`
	Category       string            `bson:"category" json:"category"`
	ImageURL       string            `bson:"image_url" json:"image_url"`
}

// Enrollment representa la estructura de una inscripción
type Enrollment struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`
	UserID  int                `json:"user_id"`
	Date    time.Time          `bson:"date" json:"date"`
}

// CreateCourse - Crear un nuevo curso (solo admin)
func CreateCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if course.Title == "" || course.Description == "" || course.Instructor == "" || 
	   course.Duration <= 0 || course.AvailableSeats <= 0 || course.Category == "" || 
	   course.ImageURL == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "All fields are required (title, description, instructor, duration, available_seats, category, image_url) and numeric values must be greater than 0"})
		return
	}

	// Validate category is one of the allowed values
	validCategories := []string{"web-development", "mobile-development", "data-science", "design", "business"}
	isValidCategory := false
	for _, validCat := range validCategories {
		if course.Category == validCat {
			isValidCategory = true
			break
		}
	}

	if !isValidCategory {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid category. Must be one of: web-development, mobile-development, data-science, design, business"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.TODO()

	// Check for duplicates
	filter := bson.M{"title": course.Title, "instructor": course.Instructor}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error checking duplicates: " + err.Error()})
		return
	}
	if count > 0 {
		jsonResponse(w, http.StatusConflict, map[string]string{"error": "A course with the same title and instructor already exists"})
		return
	}

	// Insert the course
	result, err := collection.InsertOne(ctx, course)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating course"})
		return
	}

	course.ID = result.InsertedID.(primitive.ObjectID)
	jsonResponse(w, http.StatusCreated, course)

	// Publish the course to RabbitMQ
	publishToRabbitMQ(course)
}

// GetAllCourses - Obtener todos los cursos
func GetAllCourses(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	collection := client.Database("coursesdb").Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test the connection before proceeding
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
		return
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding courses: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error finding courses"})
		return
	}
	defer cursor.Close(ctx)

	var courses []Course
	if err = cursor.All(ctx, &courses); err != nil {
		log.Printf("Error decoding courses: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error decoding courses"})
		return
	}

	if courses == nil {
		courses = []Course{} // Return empty array instead of null
	}

	jsonResponse(w, http.StatusOK, courses)
}

// GetCourse - Obtener un curso por ID
func GetCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.Background()

	var course Course
	err = collection.FindOne(ctx, bson.M{"_id": courseID}).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Course not found"})
			return
		}
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error finding course"})
		return
	}

	jsonResponse(w, http.StatusOK, course)
}

// UpdateCourse - Actualizar un curso por ID (solo admin)
func UpdateCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate all required fields
	if course.Title == "" || course.Description == "" || course.Instructor == "" || 
	   course.Duration <= 0 || course.AvailableSeats <= 0 || course.Category == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "All fields are required (title, description, instructor, duration, available_seats, category) and numeric values must be greater than 0"})
		return
	}

	// Validate category
	validCategories := []string{"web-development", "mobile-development", "data-science", "design", "business"}
	isValidCategory := false
	for _, validCat := range validCategories {
		if course.Category == validCat {
			isValidCategory = true
			break
		}
	}

	if !isValidCategory {
		jsonResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid category. Must be one of: web-development, mobile-development, data-science, design, business"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.TODO()

	// Set the ID to ensure it's preserved in the update
	course.ID = courseID

	filter := bson.M{"_id": courseID}
	update := bson.M{"$set": course}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating course"})
		return
	}

	if result.MatchedCount == 0 {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Course not found"})
		return
	}

	jsonResponse(w, http.StatusOK, course)

	// Publish the updated course to RabbitMQ
	publishToRabbitMQ(course)
}

// DeleteCourse - Eliminar un curso por ID (solo admin)
func DeleteCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.TODO()

	// Check if course exists in MongoDB
	var course Course
	err = collection.FindOne(ctx, bson.M{"_id": courseID}).Decode(&course)
	mongoExists := err != mongo.ErrNoDocuments

	// Check if course exists in Solr
	solrExists := checkSolrCourse(courseID.Hex())

	// If course doesn't exist in either system
	if !mongoExists && !solrExists {
		jsonResponse(w, http.StatusNotFound, map[string]string{"message": "Course not found"})
		return
	}

	// If course exists in MongoDB, delete it
	if mongoExists {
		_, err = collection.DeleteOne(ctx, bson.M{"_id": courseID})
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting course from MongoDB"})
			return
		}
	}

	// If course exists in Solr or was in MongoDB, delete from Solr
	if solrExists || mongoExists {
		publishDeleteEvent(courseID)
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "Course deleted successfully"})
}

// jsonResponse is a helper function to send JSON responses
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func verifyJWT(r *http.Request) (Claims, error) {
	// Get auth cookie
	cookie, err := r.Cookie("auth")
	if err != nil {
		return Claims{}, errors.New("auth cookie not found")
	}

	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return *claims, nil
	}
	return Claims{}, err
}

func CreateEnrollment(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT
	claims, err := verifyJWT(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
		return
	}

	var enrollment Enrollment
	err = json.NewDecoder(r.Body).Decode(&enrollment)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set the UserID from the JWT claims
	enrollment.UserID = claims.UserID

	// log the user id
	log.Println("UserID:", enrollment.UserID)
	log.Println("CourseID:", enrollment.CourseID)

	collection := client.Database("coursesdb").Collection("enrollments")
	ctx := context.TODO()

	// Check if the user is already enrolled in the course
	filter := bson.M{"course_id": enrollment.CourseID, "userid": enrollment.UserID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error checking enrollment: " + err.Error()})
		return
	}
	if count > 0 {
		jsonResponse(w, http.StatusConflict, map[string]string{"error": "User is already enrolled in this course"})
		return
	}

	// Check for available seats
	courseCollection := client.Database("coursesdb").Collection("courses")
	var course Course
	err = courseCollection.FindOne(ctx, bson.M{"_id": enrollment.CourseID}).Decode(&course)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error fetching course: " + err.Error()})
		return
	}
	if course.AvailableSeats <= 0 {
		jsonResponse(w, http.StatusConflict, map[string]string{"error": "No seats available for this course"})
		return
	}

	// Insert the enrollment
	result, err := collection.InsertOne(ctx, enrollment)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating enrollment"})
		return
	}

	// Decrease available seats
	_, err = courseCollection.UpdateOne(ctx, bson.M{"_id": enrollment.CourseID}, bson.M{"$inc": bson.M{"available_seats": -1}})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating available seats"})
		return
	}

	// Publish the enrollment update to RabbitMQ
	publishEnrollmentUpdateToRabbitMQ(enrollment.CourseID, -1)

	enrollment.ID = result.InsertedID.(primitive.ObjectID)
	jsonResponse(w, http.StatusCreated, enrollment)
}

func publishEnrollmentUpdateToRabbitMQ(courseID primitive.ObjectID, seatChange int) {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	update := map[string]interface{}{
		"course_id":  courseID.Hex(),
		"seat_change": seatChange,
	}

	body, err := json.Marshal(update)
	if err != nil {
		log.Printf("Failed to marshal update: %v", err)
		return
	}

	err = ch.Publish(
		"",                // exchange
		"course_updates",  // routing key
		false,             // mandatory
		false,             // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
	}
}

func CalculateAvailability(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	var courseIDs []primitive.ObjectID
	err := json.NewDecoder(r.Body).Decode(&courseIDs)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.TODO()

	var wg sync.WaitGroup
	availability := make(map[primitive.ObjectID]int)
	mu := sync.Mutex{}

	// TODO: Implement canals.
	for _, courseID := range courseIDs {
		wg.Add(1)
		go func(id primitive.ObjectID) {
			defer wg.Done()
			var course Course
			err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&course)
			if err != nil {
				return
			}
			mu.Lock()
			availability[id] = course.AvailableSeats
			mu.Unlock()
		}(courseID)
	}

	wg.Wait()
	jsonResponse(w, http.StatusOK, availability)
}

func CheckEnrollment(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	claims, err := verifyJWT(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
		return
	}

	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["course_id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	collection := client.Database("coursesdb").Collection("enrollments")
	ctx := context.TODO()

	filter := bson.M{"course_id": courseID, "userid": claims.UserID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error checking enrollment: " + err.Error()})
		return
	}

	if count > 0 {
		jsonResponse(w, http.StatusOK, map[string]string{"message": "User is enrolled in this course"})
	} else {
		jsonResponse(w, http.StatusOK, map[string]string{"message": "User is not enrolled in this course"})
	}
}

func publishToRabbitMQ(course Course) {
	rabbitmqURI := os.Getenv("RABBITMQ_URI")
	if rabbitmqURI == "" {
		rabbitmqURI = "amqp://guest:guest@rabbitmq:5672/"
	}

	var conn *amqp091.Connection
	var err error

	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		conn, err = amqp091.Dial(rabbitmqURI)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Printf("Failed to connect to RabbitMQ after 5 attempts: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
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
		log.Printf("Failed to declare queue: %v", err)
		return
	}

	// Create event payload with the correct course ID
	event := map[string]interface{}{
		"action": "upsert",
		"course": map[string]interface{}{
			"id":              course.ID.Hex(), // Convert ObjectID to string
			"title":           course.Title,
			"description":     course.Description,
			"category":        course.Category,
			"instructor":      course.Instructor,
			"duration":        course.Duration,
			"available_seats": course.AvailableSeats,
			"image_url":       course.ImageURL,
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	// Add debug logging
	log.Printf("Publishing course update to RabbitMQ: %s", string(body))

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
			amqp091.Publishing{
				DeliveryMode: amqp091.Persistent,
				ContentType:  "application/json",
				Body:        body,
			})
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return
	}

	log.Printf("Successfully published course update for ID: %s", course.ID.Hex())
}

func CheckAvailability(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	var courseIDs []string
	if err := json.NewDecoder(r.Body).Decode(&courseIDs); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	var objectIDs []primitive.ObjectID
	for _, id := range courseIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid course ID format"})
			return
		}
		objectIDs = append(objectIDs, objectID)
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.Background()

	filter := bson.M{
		"_id": bson.M{"$in": objectIDs},
		"available_seats": bson.M{"$gt": 0},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error checking availability"})
		return
	}
	defer cursor.Close(ctx)

	var availableCourses []Course
	if err = cursor.All(ctx, &availableCourses); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error processing courses"})
		return
	}

	jsonResponse(w, http.StatusOK, availableCourses)
}

func publishDeleteEvent(courseID primitive.ObjectID) {
	// Connect to RabbitMQ
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %v", err)
		return
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
		log.Printf("Failed to declare queue: %v", err)
		return
	}

	event := map[string]interface{}{
		"action": "delete",
		"course": map[string]interface{}{
			"id": courseID.Hex(),
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
}

func checkSolrCourse(courseID string) bool {
	solrURL := fmt.Sprintf("http://solr:8983/solr/courses/select?q=id:%s&wt=json", url.QueryEscape(courseID))
	resp, err := http.Get(solrURL)
	if err != nil {
		log.Printf("Error querying SolR: %v", err)
		return false
	}
	defer resp.Body.Close()

	var solrResponse struct {
		Response struct {
			NumFound int `json:"numFound"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&solrResponse); err != nil {
		log.Printf("Error decoding SolR response: %v", err)
		return false
	}

	return solrResponse.Response.NumFound > 0
}

// GetUserCourses retrieves all courses that a user is enrolled in
func GetUserCourses(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL parameters
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Verify JWT token
	claims, err := verifyJWT(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
		return
	}

	// Ensure the requesting user matches the user_id in the URL
	if claims.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]string{"error": "Forbidden: Cannot access other user's courses"})
		return
	}

	// First get all enrollments for the user
	enrollmentsCollection := client.Database("coursesdb").Collection("enrollments")
	coursesCollection := client.Database("coursesdb").Collection("courses")
	ctx := context.Background()

	// Find all enrollments for the user
	cursor, err := enrollmentsCollection.Find(ctx, bson.M{"userid": userID})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error finding enrollments"})
		return
	}
	defer cursor.Close(ctx)

	var enrollments []Enrollment
	if err = cursor.All(ctx, &enrollments); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error decoding enrollments"})
		return
	}

	// Get all course IDs from enrollments
	var courseIDs []primitive.ObjectID
	for _, enrollment := range enrollments {
		courseIDs = append(courseIDs, enrollment.CourseID)
	}

	// If user has no enrollments, return empty array
	if len(courseIDs) == 0 {
		jsonResponse(w, http.StatusOK, map[string]interface{}{"results": []Course{}})
		return
	}

	// Find all courses that match the course IDs
	cursor, err = coursesCollection.Find(ctx, bson.M{"_id": bson.M{"$in": courseIDs}})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error finding courses"})
		return
	}
	defer cursor.Close(ctx)

	var courses []Course
	if err = cursor.All(ctx, &courses); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error decoding courses"})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"results": courses})
}
