package users

import "github.com/golang-jwt/jwt/v4"

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
	Username         string `json:"username"`
	UserID           int    `json:"user_id"`
	Admin            bool   `json:"admin"`
	RegisteredClaims jwt.RegisteredClaims
}

// Implement the Valid method for the Claims struct
func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
