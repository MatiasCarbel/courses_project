package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"users-api/models"
	"users-api/repositories"
	"users-api/services"

	"github.com/stretchr/testify/assert"
)

func setupTestHandlers() *UserHandlers {
	mockRepo := repositories.NewMockUserRepository()
	service := services.NewUserService(mockRepo, []byte("test-secret"))
	return NewUserHandlers(service)
}

func TestCreateUser(t *testing.T) {
	handlers := setupTestHandlers()

	t.Run("successful user creation", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handlers.CreateUser(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEqual(t, 0, response.ID)
		assert.Equal(t, user.Username, response.Username)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()

		handlers.CreateUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginUser(t *testing.T) {
	handlers := setupTestHandlers()

	// First create a user
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	handlers.service.CreateUser(&user)

	t.Run("successful login", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handlers.LoginUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response.Token)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "wrong@example.com",
			Password: "wrongpassword",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handlers.LoginUser(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
} 