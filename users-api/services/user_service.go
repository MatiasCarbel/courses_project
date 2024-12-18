package services

import (
	"errors"
	"time"
	"users-api/models"
	"users-api/repositories"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      repositories.UserRepository
	jwtSecret []byte
}

func NewUserService(repo repositories.UserRepository, jwtSecret []byte) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("all fields are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = s.repo.Create(user)
	if err != nil {
		if err.Error() == "email already exists" {
			return err
		}
		return errors.New("error creating user")
	}

	user.Password = "" // Clear password before returning
	return nil
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user.Username, user.ID, user.Admin)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) generateJWT(username string, userID int, admin bool) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"user_id":  userID,
		"admin":    admin,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
} 