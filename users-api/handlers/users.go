package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("uccdemy") // Reemplazar por una clave segura en producción
var mc = memcache.New("localhost:11211")  // Cliente Memcached

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

// Estructura de Claims para el JWT
type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	Admin    bool   `json:"admin"`
	RegisteredClaims jwt.RegisteredClaims
}

// Implement the Valid method for the Claims struct
func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}

// Función auxiliar para generar JWT con flag admin
func generateJWT(username string, userID int, admin bool) (string, error) {
	claims := Claims{
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

// CreateUser - Crear un nuevo usuario y devolverlo sin la contraseña
func CreateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
			jsonResponse(w, http.StatusBadRequest, "Cuerpo de la solicitud inválido", "")
			return
	}

	// Debug: Imprimir los valores recibidos
	log.Printf("Datos recibidos: Username=%s, Email=%s, Password=%s", user.Username, user.Email, user.Password)

	// Verificar que todos los campos requeridos estén presentes
	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
			jsonResponse(w, http.StatusBadRequest, "Todos los campos son obligatorios", "")
			return
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
			log.Printf("Error al hashear la contraseña: %v", err)
			jsonResponse(w, http.StatusInternalServerError, "Error al crear el usuario", "")
			return
	}

	// Insertar el usuario en la base de datos
	query := "INSERT INTO users (username, email, password_hash, admin) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, user.Username, user.Email, string(hashedPassword), user.Admin)
	if err != nil {
			log.Printf("Error al crear el usuario: %v", err)
			jsonResponse(w, http.StatusConflict, "El nombre de usuario o correo electrónico ya existe", "")
			return
	}

	// Obtener el ID del usuario recién creado
	userID, err := result.LastInsertId()
	if err != nil {
			log.Printf("Error al obtener el ID del nuevo usuario: %v", err)
			jsonResponse(w, http.StatusInternalServerError, "Error al crear el usuario", "")
			return
	}

	// Devolver los datos del usuario recién creado (sin la contraseña)
	createdUser := User{
			ID:       int(userID),
			Username: user.Username,
			Email:    user.Email,
			Admin:    user.Admin, // Incluir el flag admin
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUser - Obtener un usuario por ID, sin devolver la contraseña
func GetUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	cacheKey := "user_" + strconv.Itoa(userID)
	cachedUser, err := mc.Get(cacheKey)
	var user User

	if err == nil {
		_ = json.Unmarshal(cachedUser.Value, &user)
		log.Println("Cache hit")
	} else {
		log.Println("Cache miss")
		query := "SELECT id, username, email, admin FROM users WHERE id = ?"
		err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Admin)
		if err != nil {
			jsonResponse(w, http.StatusNotFound, "Usuario no encontrado", "")
			return
		}
		userData, _ := json.Marshal(user)
		mc.Set(&memcache.Item{Key: cacheKey, Value: userData, Expiration: int32(300)})
	}

	// Devolver solo ID, username, email y el flag admin, sin la contraseña
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser - Actualizar un usuario (autenticado)
func UpdateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDParam, _ := strconv.Atoi(vars["id"])

	claims, err := verifyJWT(r)
	if err != nil || claims.UserID != userIDParam {
		jsonResponse(w, http.StatusUnauthorized, "No autorizado para actualizar este usuario", "")
		return
	}

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	query := "UPDATE users SET username = ?, email = ? WHERE id = ?"
	_, err = db.Exec(query, user.Username, user.Email, claims.UserID)
	if err != nil {
		log.Printf("Error actualizando usuario: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "Error al actualizar el usuario", "")
		return
	}

	mc.Delete("user_" + strconv.Itoa(claims.UserID)) // Invalidate cache
	jsonResponse(w, http.StatusOK, "Usuario actualizado exitosamente", "")
}

// DeleteUser - Eliminar un usuario solo si es el mismo que está autenticado y existe
func DeleteUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDParam, _ := strconv.Atoi(vars["id"])

	claims, err := verifyJWT(r)
	if err != nil || claims.UserID != userIDParam {
		jsonResponse(w, http.StatusUnauthorized, "No autorizado para eliminar este usuario", "")
		return
	}

	// Verificar si el usuario existe antes de eliminarlo
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"
	err = db.QueryRow(checkQuery, claims.UserID).Scan(&exists)
	if err != nil || !exists {
		jsonResponse(w, http.StatusNotFound, "Usuario no encontrado o ya eliminado", "")
		return
	}

	// Proceder con la eliminación si el usuario existe
	query := "DELETE FROM users WHERE id = ?"
	_, err = db.Exec(query, claims.UserID)
	if err != nil {
		log.Printf("Error eliminando usuario: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "Error al eliminar el usuario", "")
		return
	}

	// Invalidate cache
	mc.Delete("user_" + strconv.Itoa(claims.UserID))

	jsonResponse(w, http.StatusOK, "Usuario eliminado exitosamente", "")
}

// LoginUser - Login de usuario y generación de JWT
func LoginUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	var user User
	var storedHash string
	query := "SELECT id, username, password_hash, admin FROM users WHERE email = ?"
	err = db.QueryRow(query, input.Email).Scan(&user.ID, &user.Username, &storedHash, &user.Admin)
	if err != nil {
			if err == sql.ErrNoRows {
					jsonResponse(w, http.StatusUnauthorized, "Usuario o contraseña incorrectos", "")
			} else {
					log.Printf("Error querying database: %v", err)
					jsonResponse(w, http.StatusInternalServerError, "Error interno del servidor", "")
			}
			return
	}

	// Comparar el hash de la contraseña almacenada con la ingresada
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(input.Password))
	if err != nil {
			log.Printf("Password comparison failed: %v", err)
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

// jsonResponse - Enviar una respuesta en formato JSON
func jsonResponse(w http.ResponseWriter, status int, message string, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
		Token:   token,
	})
}
