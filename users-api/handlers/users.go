package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	users "users-api/domain"
	"users-api/services"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var jwtSecret = []byte("uccdemy")        // Reemplazar por una clave segura en producción
var mc = memcache.New("localhost:11211") // Cliente Memcached

// Función auxiliar para generar JWT con flag admin
func generateJWT(username string, userID int, admin bool) (string, error) {
	claims := users.Claims{
		Username: username,
		UserID:   userID,
		Admin:    admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // El token expira en 72 horas
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Función para verificar el token JWT y extraer el userID y el flag admin
func verifyJWT(r *http.Request) (users.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return users.Claims{}, errors.New("token no provisto")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &users.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return jwtSecret, nil
	})
	if claims, ok := token.Claims.(*users.Claims); ok && token.Valid {
		return *claims, nil
	}
	return users.Claims{}, err
}

type UserHandler struct {
	UserService *services.UserService
}

// GetUserHandler handles HTTP requests for getting a user
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	cacheKey := "user_" + strconv.Itoa(userID)
	cachedUser, err := mc.Get(cacheKey)
	var user users.User

	if err == nil {
		_ = json.Unmarshal(cachedUser.Value, &user)
		log.Println("Cache hit")
	} else {
		log.Println("Cache miss")
		user, err := h.UserService.GetUser(userID)
		if err != nil {
			jsonResponse(w, http.StatusNotFound, "Usuario no encontrado", "")
			return
		}
		userData, _ := json.Marshal(user)
		mc.Set(&memcache.Item{Key: cacheKey, Value: userData, Expiration: int32(300)})

		// Respond with the created user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}

}

// CreateUserHandler handles HTTP requests for creating a user
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req services.CreateUserRequest

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to create a user
	user, err := h.UserService.CreateUser(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser - Actualizar un usuario (autenticado)
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDParam, _ := strconv.Atoi(vars["id"])

	claims, err := verifyJWT(r)
	if err != nil || claims.UserID != userIDParam {
		jsonResponse(w, http.StatusUnauthorized, "No autorizado para actualizar este usuario", "")
		return
	}

	var user users.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	err = h.UserService.UpdateUser(&user, userIDParam)
	if err != nil {
		log.Printf("Error actualizando usuario: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "Error al actualizar el usuario", "")
		return
	}

	mc.Delete("user_" + strconv.Itoa(claims.UserID)) // Invalidate cache
	jsonResponse(w, http.StatusOK, "Usuario actualizado exitosamente", "")
}

// DeleteUserHandler handles HTTP requests for deleting a user
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDParam, _ := strconv.Atoi(vars["id"])

	claims, err := verifyJWT(r)
	if err != nil || claims.UserID != userIDParam {
		jsonResponse(w, http.StatusUnauthorized, "No autorizado para eliminar este usuario", "")
		return
	}

	// Verificar si el usuario existe antes de eliminarlo
	exists, err := h.UserService.UserExists(userIDParam)
	if err != nil || !exists {
		jsonResponse(w, http.StatusNotFound, "Usuario no encontrado o ya eliminado", "")
		return
	}

	// Proceder con la eliminación si el usuario existe
	err = h.UserService.DeleteUser(userIDParam)
	if err != nil {
		log.Printf("Error eliminando usuario: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "Error al eliminar el usuario", "")
		return
	}

	// Invalidate cache
	mc.Delete("user_" + strconv.Itoa(claims.UserID))

	jsonResponse(w, http.StatusOK, "Usuario eliminado exitosamente", "")

	// Respond with eliminated user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Usuario eliminado exitosamente")

}

// jsonResponse - Enviar una respuesta en formato JSON
func jsonResponse(w http.ResponseWriter, status int, message string, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(users.Response{
		Status:  status,
		Message: message,
		Token:   token,
	})
}

// LoginUser - Login de usuario y generación de JWT
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Invalid request body", "")
		return
	}

	// Verificar si los campos están presentes
	if input.Email == "" || input.Password == "" {
		jsonResponse(w, http.StatusBadRequest, "El email y la contraseña son obligatorios", "")
		return
	}

	user, err := h.UserService.LoginUser(input.Email, input.Password)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "Usuario o contraseña incorrectos", "")
		return
	}

	// Generar el JWT si la contraseña es correcta
	token, err := generateJWT(user.Username, user.ID, user.Admin)
	if err != nil {
		log.Printf("Error generando el token: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "Error generando el token", "")
		return
	}

	jsonResponse(w, http.StatusOK, "Login exitoso", token)
}
