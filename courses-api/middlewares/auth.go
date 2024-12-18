package middlewares

import (
	"context"
	"courses-api/models"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}
	if c.UserID == 0 {
		return models.ErrInvalidToken
	}
	return nil
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func VerifyAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, models.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, models.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		if !claims.Admin {
			http.Error(w, models.ErrForbidden.Error(), http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "admin", claims.Admin)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func VerifyToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// First check Authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// If no Authorization header, check cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, models.ErrUnauthorized.Error(), http.StatusUnauthorized)
				return
			}
			tokenString = cookie.Value
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, models.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
