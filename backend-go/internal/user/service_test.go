package user_test

import (
	"context"
	"database/sql"
	"task-manager/backend-go/internal/i18n"
	"task-manager/backend-go/internal/user"
	"task-manager/backend-go/models"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

const (
	testUsername = "thorbar"
	testPassword = "abcd1"
	testName     = "Thor"
	testSurname  = "Odinson"
	testEmail    = "thorbar@example.com"
)

// setupTestDB creates an in-memory SQLite database with users table for testing.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	schema := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		surname TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(schema)
	assert.NoError(t, err)

	return db
}

// TestRegisterUser verifies user registration stores data correctly.
func TestRegisterUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := user.NewService(db)

	req := &models.RegisterRequest{
		Name:     testName,
		Surname:  testSurname,
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	}

	err := service.RegisterUser(context.Background(), req)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

// TestLoginUser_Success tests successful login returns a token.
func TestLoginUser_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := user.NewService(db)

	// Register user first
	err := service.RegisterUser(context.Background(), &models.RegisterRequest{
		Name:     testName,
		Surname:  testSurname,
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)

	// Attempt login
	token, err := service.LoginUser(context.Background(), &models.LoginRequest{
		Username: testUsername,
		Password: testPassword,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

// TestLoginUser_WrongPassword ensures login fails with incorrect password.
func TestLoginUser_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := user.NewService(db)

	// Register user
	err := service.RegisterUser(context.Background(), &models.RegisterRequest{
		Name:     testName,
		Surname:  testSurname,
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)

	// Login with wrong password
	_, err = service.LoginUser(context.Background(), &models.LoginRequest{
		Username: testUsername,
		Password: "wrongpass",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), i18n.T("incorrect_password")) // Use i18n key for error message
}

// TestLoginUser_UserNotFound ensures login fails when user doesn't exist.
func TestLoginUser_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := user.NewService(db)

	// Attempt login without registration
	_, err := service.LoginUser(context.Background(), &models.LoginRequest{
		Username: "noexistente",
		Password: "any",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), i18n.T("user_not_found")) // Use i18n key for error message
}

// TestRegisterUser_DuplicateUsername checks registration fails for duplicate username.
func TestRegisterUser_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := user.NewService(db)

	req := &models.RegisterRequest{
		Name:     testName,
		Surname:  testSurname,
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	}

	err := service.RegisterUser(context.Background(), req)
	assert.NoError(t, err)

	// Change email, try duplicate username
	req.Email = "other@example.com"
	err = service.RegisterUser(context.Background(), req)
	assert.Error(t, err) // Expected to fail due to duplicate username
}

// TestRegisterUser_DuplicateEmail checks registration fails for duplicate email.
func TestRegisterUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := user.NewService(db)

	req := &models.RegisterRequest{
		Name:     testName,
		Surname:  testSurname,
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	}

	err := service.RegisterUser(context.Background(), req)
	assert.NoError(t, err)

	// Change username, try duplicate email
	req.Username = "loki_dos"
	err = service.RegisterUser(context.Background(), req)
	assert.Error(t, err) // Expected to fail due to duplicate email
}

// TestUpdateUser_Success verifies updating user fields works correctly.
func TestUpdateUser_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := user.NewService(db)

	// Step 1: Register user
	err := service.RegisterUser(context.Background(), &models.RegisterRequest{
		Name:     testName,
		Surname:  testSurname,
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.NoError(t, err)

	// Step 2: Get user ID
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", testUsername).Scan(&userID)
	assert.NoError(t, err)

	// Step 3: Update user info
	newUsername := "thorgod"
	newName := "tango"
	newSurname := "express"

	err = service.UpdateUser(userID, newName, newSurname)
	assert.NoError(t, err)

	// Step 4: Confirm DB changes
	var updatedUsername, updatedName, updatedSurname string
	err = db.QueryRow("SELECT username, name, surname FROM users WHERE id = ?", userID).
		Scan(&updatedUsername, &updatedName, &updatedSurname)
	assert.NoError(t, err)
	assert.Equal(t, newUsername, updatedUsername)
	assert.Equal(t, newName, updatedName)
	assert.Equal(t, newSurname, updatedSurname)
}
