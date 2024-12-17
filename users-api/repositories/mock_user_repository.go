package repositories

import (
	"errors"
	"users-api/models"
)

type MockUserRepository struct {
	users map[int]*models.User
	emails map[string]*models.User
	lastID int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[int]*models.User),
		emails: make(map[string]*models.User),
		lastID: 0,
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("missing required fields")
	}

	// Check for duplicate email
	for _, existingUser := range m.users {
		if existingUser.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	// Create a copy of the user to store
	storedUser := &models.User{
		ID:       m.lastID + 1,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password, // Store the hashed password
		Admin:    user.Admin,
	}

	m.lastID++
	user.ID = m.lastID
	m.users[storedUser.ID] = storedUser
	m.emails[storedUser.Email] = storedUser
	return nil
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	// Return a copy of the user without the password
	return &models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Admin:    user.Admin,
		Password: "", // Don't return the password
	}, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	user, exists := m.emails[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	// Return a copy with the password (needed for login)
	return &models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Admin:    user.Admin,
		Password: user.Password, // Include password for login verification
	}, nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(id int) error {
	user, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}
	delete(m.emails, user.Email)
	delete(m.users, id)
	return nil
} 