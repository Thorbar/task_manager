package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"task-manager/backend-go/db"
	"task-manager/backend-go/internal/auth"
	"task-manager/backend-go/internal/i18n"
	"task-manager/backend-go/models"
)

// UserRouter routes requests to the appropriate handler based on HTTP method.
func UserRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetUserHandler(w, r)
	case http.MethodPost:
		UpdateUserHandler(w, r)
	default:
		http.Error(w, i18n.T("http.error.method_not_allowed"), http.StatusMethodNotAllowed)
	}
}

// GetUserHandler retrieves the user info if the token is valid.
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserSummary

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "auth.error.token_not_provided")
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := auth.ParseToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "auth.error.invalid_token")
		return
	}

	err = db.DB.QueryRow("SELECT username, email, name, surname, created_at FROM users WHERE id = ?", userID).
		Scan(&user.Username, &user.Email, &user.Name, &user.Surname, &user.Created)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user.error.not_found")
		return
	}


	respondWithJSON(w, http.StatusOK, user)
}

// Handler groups dependencies for handlers (optional).
type Handler struct {
	Service *Service
}

// UpdateUserHandler updates user's name and surname after validating the token and request data.
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateUserRequest

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "http.error.method_not_allowed")
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "auth.error.token_not_provided")
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	userID, err := auth.ParseToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "auth.error.invalid_token")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "http.error.invalid_data")
		return
	}

	service := NewService(db.DB)
	if err := service.UpdateUser(userID, req.Name, req.Surname); err != nil {
		respondWithError(w, http.StatusInternalServerError, "user.error.update_failed")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": i18n.T("user.success.updated")})
}

// respondWithError writes a JSON error response with localized message.
func respondWithError(w http.ResponseWriter, statusCode int, msgKey string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": i18n.T(msgKey),
	})
}

// respondWithJSON writes a JSON response with data.
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
