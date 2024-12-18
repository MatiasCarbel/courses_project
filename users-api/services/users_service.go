package services

import (
	"errors"
	users "users-api/domain"
	dao "users-api/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserDAO *dao.UserDAO
}

// CreateUserRequest represents the input data required to create a user
type CreateUserRequest struct {
	Username string
	Email    string
	Password string
	Admin    bool
}

// CreateUser handles the business logic for creating a user
func (s *UserService) CreateUser(req *CreateUserRequest) (*users.User, error) {
	// Input validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Prepare the user object
	user := &users.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Admin:    req.Admin,
	}

	// Insert the user into the database
	userID, err := s.UserDAO.CreateUser(user)
	if err != nil {
		return nil, err // Propagate DAO error
	}

	// Populate the user ID
	user.ID = userID

	// Remove the password before returning (security measure)
	user.Password = ""

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID int) (*users.User, error) {
	user, err := s.UserDAO.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Return if a user exists by ID
func (s *UserService) UserExists(userID int) (bool, error) {
	exists, err := s.UserDAO.UserExists(userID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(userID int) error {
	err := s.UserDAO.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user by ID
func (s *UserService) UpdateUser(user *users.User, userID int) error {
	err := s.UserDAO.UpdateUser(user, userID)
	if err != nil {
		return err
	}

	return nil
}

// LoginUser handles the business logic for logging in a user
func (s *UserService) LoginUser(email, password string) (*users.User, error) {
	// Retrieve the user by email
	user, err := s.UserDAO.LoginUser(email)
	if err != nil {
		return nil, err
	}

	// Compare the password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// Remove the password before returning (security measure)
	user.Password = ""

	return user, nil
}
