package services

import (
	"testing"
	"users-api/models"
	"users-api/repositories"

	"github.com/stretchr/testify/assert"
)

func setupTestService() (*UserService, *repositories.MockUserRepository) {
	mockRepo := repositories.NewMockUserRepository()
	service := NewUserService(mockRepo, []byte("test-secret"))
	return service, mockRepo
}

func TestGetUserByID(t *testing.T) {
	service, _ := setupTestService()

	// Create a test user first
	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Admin:    false,
	}
	err := service.CreateUser(testUser)
	assert.NoError(t, err)
	userID := testUser.ID

	t.Run("successful user retrieval", func(t *testing.T) {
		user, err := service.GetUserByID(userID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.Username, user.Username)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Admin, user.Admin)
		assert.Empty(t, user.Password, "Password should be empty when retrieving user")
	})

	t.Run("user not found", func(t *testing.T) {
		user, err := service.GetUserByID(999)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestCreateUser(t *testing.T) {
	service, _ := setupTestService()

	t.Run("successful user creation", func(t *testing.T) {
		user := &models.User{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			Admin:    false,
		}

		err := service.CreateUser(user)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, user.ID)
		assert.Empty(t, user.Password) // Password should be cleared
	})

	t.Run("missing required fields", func(t *testing.T) {
		user := &models.User{
			Username: "",
			Email:    "test@example.com",
			Password: "password123",
		}

		err := service.CreateUser(user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "all fields are required")
	})

	t.Run("duplicate email", func(t *testing.T) {
		// Create first user
		user1 := &models.User{
			Username: "user1",
			Email:    "duplicate@example.com",
			Password: "password123",
			Admin:    false,
		}
		err := service.CreateUser(user1)
		assert.NoError(t, err)

		// Try to create second user with same email
		user2 := &models.User{
			Username: "user2",
			Email:    "duplicate@example.com",
			Password: "password456",
			Admin:    false,
		}
		err = service.CreateUser(user2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already exists")
	})
}

func TestLogin(t *testing.T) {
	service, _ := setupTestService()

	// Create a test user first
	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Admin:    false,
	}
	err := service.CreateUser(testUser)
	assert.NoError(t, err)

	// Store the email for login test
	userEmail := testUser.Email

	t.Run("successful login", func(t *testing.T) {
		token, err := service.Login(userEmail, "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("invalid email", func(t *testing.T) {
		token, err := service.Login("wrong@example.com", "password123")
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "invalid credentials")
	})

	t.Run("invalid password", func(t *testing.T) {
		token, err := service.Login(userEmail, "wrongpassword")
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "invalid credentials")
	})
}

func TestGenerateJWT(t *testing.T) {
	service, _ := setupTestService()

	t.Run("successful token generation", func(t *testing.T) {
		token, err := service.generateJWT("testuser", 1, false)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
