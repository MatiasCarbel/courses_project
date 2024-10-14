package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Course struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

// Datos de ejemplo
var courses = []Course{
    {ID: 1, Name: "Curso de Go", Description: "Aprende Go desde cero"},
    {ID: 2, Name: "Curso de Docker", Description: "Domina Docker y contenedores"},
}

// Obtener todos los cursos
func GetCourses(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(courses)
}

// Obtener curso por ID
func GetCourseByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    // Búsqueda de cursos (esto luego se puede optimizar con base de datos)
    for _, item := range courses {
        if fmt.Sprintf("%d", item.ID) == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    http.Error(w, "Curso no encontrado", http.StatusNotFound)
}

// Crear un nuevo curso
func CreateCourse(w http.ResponseWriter, r *http.Request) {
    var course Course
    _ = json.NewDecoder(r.Body).Decode(&course)
    course.ID = len(courses) + 1
    courses = append(courses, course)
    json.NewEncoder(w).Encode(course)
}

// Actualizar curso
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var updatedCourse Course
    _ = json.NewDecoder(r.Body).Decode(&updatedCourse)
    for i, item := range courses {
        if fmt.Sprintf("%d", item.ID) == params["id"] {
            courses[i] = updatedCourse
            json.NewEncoder(w).Encode(updatedCourse)
            return
        }
    }
    http.Error(w, "Curso no encontrado", http.StatusNotFound)
}

// Eliminar curso
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for i, item := range courses {
        if fmt.Sprintf("%d", item.ID) == params["id"] {
            courses = append(courses[:i], courses[i+1:]...)
            break
        }
    }
    w.WriteHeader(http.StatusNoContent)
}
