package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("uccdemy") // Use a secure key in production

// Claims define la estructura de los claims del JWT
type Claims struct {
	Admin bool `json:"admin"`
	jwt.RegisteredClaims
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

// Middleware para verificar si el usuario es admin
func VerifyAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// claims, err := verifyJWT(r)
		// fmt.Println("claims: ", claims)
		// fmt.Println("err: ", err)
		// if err != nil {
		// 	jsonResponse(w, http.StatusUnauthorized, "Unauthorized", "")
		// 	return
		// }

		// if !claims.Admin {
		// 	jsonResponse(w, http.StatusForbidden, "Admin access required", "")
		// 	return
		// }

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
