package tests

import (
	"database/sql"
	"os"
	"testing"
	"users-api/models"
	"users-api/repositories"
	"users-api/services"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var (
	testDB     *sql.DB
	testService *services.UserService
)

func TestMain(m *testing.M) {
	// Set up test database connection
	var err error
	testDB, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/courses_test")
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	// Create test tables
	setupTestDB()

	// Run tests
	code := m.Run()

	// Clean up
	cleanupTestDB()

	os.Exit(code)
}

func setupTestDB() {
	// Create users table
	_, err := testDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			admin BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		panic(err)
	}
}

func cleanupTestDB() {
	_, err := testDB.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		panic(err)
	}
}

func setupTest(t *testing.T) *services.UserService {
	// Clear existing data
	_, err := testDB.Exec("DELETE FROM users")
	assert.NoError(t, err)

	// Create repository and service
	repo := repositories.NewSQLUserRepository(testDB)
	return services.NewUserService(repo, []byte("test-secret"))
}

func TestIntegrationCreateUser(t *testing.T) {
	service := setupTest(t)

	t.Run("successful user creation", func(t *testing.T) {
		user := &models.User{
			Username: "integrationtest",
			Email:    "integration@test.com",
			Password: "testpass123",
		}

		err := service.CreateUser(user)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, user.ID)

		// Verify user exists in database
		var dbUser models.User
		err = testDB.QueryRow("SELECT id, username, email FROM users WHERE id = ?", user.ID).
			Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email)
		assert.NoError(t, err)
		assert.Equal(t, user.Username, dbUser.Username)
	})
}

func TestIntegrationLogin(t *testing.T) {
	service := setupTest(t)

	// Create test user
	user := &models.User{
		Username: "logintest",
		Email:    "login@test.com",
		Password: "testpass123",
	}
	err := service.CreateUser(user)
	assert.NoError(t, err)

	t.Run("successful login", func(t *testing.T) {
		token, err := service.Login("login@test.com", "testpass123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
} 