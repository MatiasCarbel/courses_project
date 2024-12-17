package repositories

import (
	"database/sql"
	"users-api/models"
)

type UserRepository interface {
    Create(user *models.User) error
    GetByID(id int) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id int) error
}

type SQLUserRepository struct {
    db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) UserRepository {
    return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) Create(user *models.User) error {
    query := "INSERT INTO users (username, email, password_hash, admin) VALUES (?, ?, ?, ?)"
    result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Admin)
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    user.ID = int(id)
    return nil
}

func (r *SQLUserRepository) GetByID(id int) (*models.User, error) {
    var user models.User
    query := "SELECT id, username, email, admin FROM users WHERE id = ?"
    err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Admin)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *SQLUserRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    var hashedPassword string
    query := "SELECT id, username, password_hash, admin FROM users WHERE email = ?"
    err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &hashedPassword, &user.Admin)
    if err != nil {
        return nil, err
    }
    user.Password = hashedPassword
    return &user, nil
}

func (r *SQLUserRepository) Update(user *models.User) error {
    query := "UPDATE users SET username = ?, email = ? WHERE id = ?"
    _, err := r.db.Exec(query, user.Username, user.Email, user.ID)
    return err
}

func (r *SQLUserRepository) Delete(id int) error {
    query := "DELETE FROM users WHERE id = ?"
    _, err := r.db.Exec(query, id)
    return err
} 