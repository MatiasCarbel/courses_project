package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
)

// Claims define la estructura de los claims del JWT
type Claims struct {
	Admin bool `json:"admin"`
	jwt.RegisteredClaims
}

// Middleware para verificar si el usuario es admin
func VerifyAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			jsonResponse(w, http.StatusUnauthorized, "Falta el token de autorización", "")
			return
		}

		// Remove "Bearer " prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("uccdemy"), nil
		})

		if err != nil || !token.Valid {
			jsonResponse(w, http.StatusUnauthorized, "Token inválido", "")
			return
		}

		// Verify if the user has admin permissions
		if !claims.Admin {
			jsonResponse(w, http.StatusForbidden, "No tienes permisos de administrador", "")
			return
		}

		next(w, r)
	}
}

// jsonResponse - Enviar una respuesta en formato JSON
func jsonResponse(w http.ResponseWriter, status int, message string, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  strconv.Itoa(status),
		"message": message,
		"token":   token,
	})
}
