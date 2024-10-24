package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
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
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Instructor  string             `bson:"instructor" json:"instructor"`
	Duration    int                `bson:"duration" json:"duration"`
	AvailableSeats int               `bson:"available_seats" json:"available_seats"`
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
	_ = json.NewDecoder(r.Body).Decode(&course)

	if course.Title == "" || course.Description == "" || course.Instructor == "" || course.Duration <= 0 || course.AvailableSeats <= 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Todos los campos son obligatorios y la duración y los asientos disponibles deben ser mayores a 0"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.TODO()

	// Check for duplicates
	filter := bson.M{"title": course.Title, "instructor": course.Instructor}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error verificando duplicados: " + err.Error()})
		return
	}
	if count > 0 {
		jsonResponse(w, http.StatusConflict, map[string]string{"error": "Ya existe un curso con el mismo título e instructor"})
		return
	}

	// Insert the course
	result, err := collection.InsertOne(ctx, course)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error al crear el curso"})
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
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error obteniendo los cursos"})
		return
	}
	defer cursor.Close(ctx)

	var courses []Course
	if err = cursor.All(ctx, &courses); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error escaneando los cursos"})
		return
	}

	jsonResponse(w, http.StatusOK, courses)
}

// GetCourse - Obtener un curso por ID
func GetCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "ID de curso inválido"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.TODO()

	var course Course
	filter := bson.M{"_id": courseID}
	err = collection.FindOne(ctx, filter).Decode(&course)
	if err == mongo.ErrNoDocuments {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Curso no encontrado"})
		return
	} else if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error obteniendo el curso"})
		return
	}

	jsonResponse(w, http.StatusOK, course)
}

// UpdateCourse - Actualizar un curso por ID (solo admin)
func UpdateCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "ID de curso inválido"})
		return
	}

	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)

	if course.Title == "" || course.Description == "" || course.Instructor == "" || course.Duration <= 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Todos los campos son obligatorios y la duración debe ser mayor a 0"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.TODO()

	filter := bson.M{"_id": courseID}
	update := bson.M{"$set": course}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error actualizando el curso"})
		return
	}

	jsonResponse(w, http.StatusOK, course)

	// Publish the course to RabbitMQ
	publishToRabbitMQ(course)
}

// DeleteCourse - Eliminar un curso por ID (solo admin)
func DeleteCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "ID de curso inválido"})
		return
	}

	collection := client.Database("coursesdb").Collection("courses")
	ctx := context.TODO()

	// Retrieve the course details before deletion
	var course Course
	err = collection.FindOne(ctx, bson.M{"_id": courseID}).Decode(&course)
	if err == mongo.ErrNoDocuments {
		jsonResponse(w, http.StatusNotFound, map[string]string{"message": "No se encontró el curso para eliminar"})
		return
	} else if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error obteniendo el curso"})
		return
	}

	// Delete the course
	result, err := collection.DeleteOne(ctx, bson.M{"_id": courseID})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error eliminando el curso"})
		return
	}

	if result.DeletedCount == 0 {
		jsonResponse(w, http.StatusNotFound, map[string]string{"message": "No se encontró el curso para eliminar"})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "Curso eliminado exitosamente"})

	// Publish the course to RabbitMQ
	publishToRabbitMQ(course)
}

// jsonResponse is a helper function to send JSON responses
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func verifyJWT(r *http.Request) (Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return Claims{}, errors.New("token no provisto")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
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

	enrollment.ID = result.InsertedID.(primitive.ObjectID)
	jsonResponse(w, http.StatusCreated, enrollment)
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
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	body, err := json.Marshal(course)
	if err != nil {
		log.Fatalf("Failed to marshal course: %v", err)
	}

	err = ch.Publish(
		"",              // exchange
		"course_updates", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}
