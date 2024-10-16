package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Course representa la estructura de un curso
type Course struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Instructor  string             `bson:"instructor" json:"instructor"`
	Duration    int                `bson:"duration" json:"duration"`
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

	if course.Title == "" || course.Description == "" || course.Instructor == "" || course.Duration <= 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Todos los campos son obligatorios y la duración debe ser mayor a 0"})
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

	filter := bson.M{"_id": courseID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error eliminando el curso"})
		return
	}

	if result.DeletedCount == 0 {
		jsonResponse(w, http.StatusNotFound, map[string]string{"message": "No se encontró el curso para eliminar"})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "Curso eliminado exitosamente"})
}

// jsonResponse is a helper function to send JSON responses
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func CreateEnrollment(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	var enrollment struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`
		UserID   int                `json:"user_id"`
		Date     time.Time          `bson:"date" json:"date"`
	}

	err := json.NewDecoder(r.Body).Decode(&enrollment)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	collection := client.Database("coursesdb").Collection("enrollments")
	ctx := context.TODO()

	result, err := collection.InsertOne(ctx, enrollment)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating enrollment"})
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
			filter := bson.M{"_id": id}
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				return
			}
			mu.Lock()
			availability[id] = int(count)
			mu.Unlock()
		}(courseID)
	}

	wg.Wait()
	jsonResponse(w, http.StatusOK, availability)
}
