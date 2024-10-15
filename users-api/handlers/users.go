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

var jwtSecret = []byte("your_secret_key") // Reemplazar por una clave segura en producción
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
    Password string `json:"-"`
}

// Función auxiliar para generar JWT
func generateJWT(username string, userID int) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "user_id":  userID,
        "exp":      time.Now().Add(72 * time.Hour).Unix(), // El token expira en 72 horas
    })
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

// Función para verificar el token JWT y extraer el userID
func verifyJWT(r *http.Request) (int, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return 0, errors.New("token no provisto")
    }
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("método de firma inválido")
        }
        return jwtSecret, nil
    })
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := int(claims["user_id"].(float64))
        return userID, nil
    }
    return 0, err
}

// CreateUser - Crear un nuevo usuario y devolverlo sin la contraseña
func CreateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
			log.Printf("Error al hashear la contraseña: %v", err)
			jsonResponse(w, http.StatusInternalServerError, "Error al crear el usuario", "")
			return
	}
	user.Password = string(hashedPassword)

	// Insertar el usuario en la base de datos
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	result, err := db.Exec(query, user.Username, user.Email, user.Password)
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
			query := "SELECT id, username, email FROM users WHERE id = ?"
			err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email)
			if err != nil {
					jsonResponse(w, http.StatusNotFound, "Usuario no encontrado", "")
					return
			}
			userData, _ := json.Marshal(user)
			mc.Set(&memcache.Item{Key: cacheKey, Value: userData, Expiration: int32(300)})
	}

	// Devolver solo ID, username y email, sin la contraseña
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser - Actualizar un usuario (autenticado)
func UpdateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDParam, _ := strconv.Atoi(vars["id"])

    userIDToken, err := verifyJWT(r)
    if err != nil || userIDToken != userIDParam {
        jsonResponse(w, http.StatusUnauthorized, "No autorizado para actualizar este usuario", "")
        return
    }

    var user User
    _ = json.NewDecoder(r.Body).Decode(&user)

    query := "UPDATE users SET username = ?, email = ? WHERE id = ?"
    _, err = db.Exec(query, user.Username, user.Email, userIDToken)
    if err != nil {
        log.Printf("Error actualizando usuario: %v", err)
        jsonResponse(w, http.StatusInternalServerError, "Error al actualizar el usuario", "")
        return
    }

    mc.Delete("user_" + strconv.Itoa(userIDToken)) // Invalidate cache
    jsonResponse(w, http.StatusOK, "Usuario actualizado exitosamente", "")
}

// DeleteUser - Eliminar un usuario solo si es el mismo que está autenticado y existe
func DeleteUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDParam, _ := strconv.Atoi(vars["id"])

	// Verificar el JWT y obtener el userID del token
	userIDToken, err := verifyJWT(r)
	if err != nil || userIDToken != userIDParam {
			jsonResponse(w, http.StatusUnauthorized, "No autorizado para eliminar este usuario", "")
			return
	}

	// Verificar si el usuario existe antes de eliminarlo
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"
	err = db.QueryRow(checkQuery, userIDToken).Scan(&exists)
	if err != nil || !exists {
			jsonResponse(w, http.StatusNotFound, "Usuario no encontrado o ya eliminado", "")
			return
	}

	// Proceder con la eliminación si el usuario existe
	query := "DELETE FROM users WHERE id = ?"
	_, err = db.Exec(query, userIDToken)
	if err != nil {
			log.Printf("Error eliminando usuario: %v", err)
			jsonResponse(w, http.StatusInternalServerError, "Error al eliminar el usuario", "")
			return
	}

	// Invalidate cache
	mc.Delete("user_" + strconv.Itoa(userIDToken))

	jsonResponse(w, http.StatusOK, "Usuario eliminado exitosamente", "")
}

// LoginUser - Login de usuario y generación de JWT
func LoginUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var input User
    _ = json.NewDecoder(r.Body).Decode(&input)

    var user User
    query := "SELECT id, username, password_hash FROM users WHERE email = ?"
    err := db.QueryRow(query, input.Email).Scan(&user.ID, &user.Username, &user.Password)
    if err == sql.ErrNoRows {
        jsonResponse(w, http.StatusUnauthorized, "Usuario o contraseña incorrectos", "")
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        jsonResponse(w, http.StatusUnauthorized, "Usuario o contraseña incorrectos", "")
        return
    }

    token, err := generateJWT(user.Username, user.ID)
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
