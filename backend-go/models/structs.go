package models

import "time"

// LoginRequest represents the payload for user login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents the payload for user registration.
type RegisterRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Task represents a user task with metadata.
type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRequest represents user data submitted from the frontend for update.
type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserSummary represents the data returned from API after user login.
type UserSummary struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Created  string `json:"created_at"`
}

// UpdateUserRequest represents the payload to update user profile data.
type UpdateUserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}
