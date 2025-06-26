package user

import (
	"context"
	"database/sql"
	"errors"
	"task-manager/backend-go/internal/auth"
	"task-manager/backend-go/internal/i18n"
	"task-manager/backend-go/models"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

// RegisterUser handles user registration logic.
// It hashes the password and inserts the user into the database.
func (s *Service) RegisterUser(ctx context.Context, req *models.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, surname, username, email, password) VALUES (?, ?, ?, ?, ?)`
	_, err = s.DB.ExecContext(ctx, query, req.Name, req.Surname, req.Username, req.Email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

// LoginUser verifies user credentials and returns a JWT token upon success.
func (s *Service) LoginUser(ctx context.Context, req *models.LoginRequest) (string, error) {
	var hashedPassword string
	var userID int

	err := s.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE username = ?", req.Username).
		Scan(&userID, &hashedPassword)

	if err == sql.ErrNoRows {
		return "", errors.New(i18n.T("user.error.not_found"))

	} else if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return "", errors.New(i18n.T("user.error.incorrect_password"))

	}

	tokenString, err := auth.GenerateJWT(userID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// UpdateUser updates user's name and surname based on user ID.
func (s *Service) UpdateUser(userID int, name, surname string) error {
	query := `UPDATE users SET name = ?, surname = ? WHERE id = ?`
	_, err := s.DB.Exec(query, name, surname, userID)
	return err
}
