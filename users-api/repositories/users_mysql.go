package dao

import (
	"database/sql"
	users "users-api/domain"
)

type UserDAO struct {
	DB *sql.DB
}

// CreateUser inserts a new user into the database.
func (dao *UserDAO) CreateUser(user *users.User) (int, error) {
	query := "INSERT INTO users (username, email, password_hash, admin) VALUES (?, ?, ?, ?)"
	result, err := dao.DB.Exec(query, user.Username, user.Email, user.Password, user.Admin)
	if err != nil {
		return 0, err
	}

	// Retrieve the inserted user's ID
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userID), nil
}

// GetUser retrieves a user from the database by ID.
func (dao *UserDAO) GetUser(userID int) (*users.User, error) {
	query := "SELECT id, username, email, admin FROM users WHERE id = ?"
	var user users.User
	err := dao.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Admin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Check if the user exists in the database
func (dao *UserDAO) UserExists(userID int) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE id = ?"
	var count int
	err := dao.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateUser updates a user in the database by ID.
func (dao *UserDAO) UpdateUser(user *users.User, paramID int) error {
	query := "UPDATE users SET username = ?, email = ? WHERE id = ?"
	_, err := dao.DB.Exec(query, user.Username, user.Email, paramID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user from the database by ID.
func (dao *UserDAO) DeleteUser(userID int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := dao.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// login
func (dao *UserDAO) LoginUser(email string) (*users.User, error) {
	query := "SELECT id, username, password_hash, admin FROM users WHERE email = ?"
	var user users.User
	err := dao.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Password, &user.Admin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
