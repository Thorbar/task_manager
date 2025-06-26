// File to create tasks
package task

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"task-manager/backend-go/db"
	"task-manager/backend-go/internal/i18n"
	"task-manager/backend-go/models"

	"github.com/golang-jwt/jwt/v5"
)

type Task struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// TasksHandler validates JWT, retrieves tasks for the user, and returns JSON response.
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, i18n.T("invalid_token_format"), http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, i18n.T("invalid_or_expired_token"), http.StatusUnauthorized)
		return
	}

	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		http.Error(w, i18n.T("missing_or_invalid_user_id"), http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, i18n.T("invalid_user_id"), http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query(`
		SELECT t.id, t.user_id, u.username, t.title, t.description, t.status, t.created_at
		FROM tasks t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id = ?`, userID)
	if err != nil {
		http.Error(w, i18n.T("error_fetching_tasks"), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Username, &t.Title, &t.Description, &t.Status, &t.CreatedAt); err != nil {
			http.Error(w, i18n.T("error_reading_tasks"), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}