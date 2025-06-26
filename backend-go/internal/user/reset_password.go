package user

import (
	"encoding/json"
	"net/http"
	"task-manager/backend-go/db"
	"task-manager/backend-go/internal/i18n"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ResetPasswordRequest represents the expected JSON payload for password reset requests.
type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ResetPasswordHandler handles HTTP requests to reset a user's password.
// It validates method, input data, token validity, password strength,
// hashes new password and updates it in the database.
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": i18n.T("method_not_allowed")})
		return
	}

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": i18n.T("invalid_request_data")})
		return
	}

	if req.Token == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": i18n.T("missing_token_or_password")})
		return
	}

	// Validate token and expiration time from DB
	var userID int
	var expirationTime time.Time
	err := db.DB.QueryRow(
		"SELECT id, password_reset_expiration FROM users WHERE password_reset_token = ?", req.Token,
	).Scan(&userID, &expirationTime)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": i18n.T("invalid_or_expired_token")})
		return
	}

	if time.Now().After(expirationTime) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": i18n.T("token_expired")})
		return
	}

	// Validate new password strength using shared utility (assumed implemented elsewhere)
	if !isValidPassword(req.Password) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("password_requirements"),
		})
		return
	}

	// Hash the new password securely
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": i18n.T("password_hash_error")})
		return
	}

	// Update user password and clear reset token/expiration in DB
	_, err = db.DB.Exec(
		"UPDATE users SET password = ?, password_reset_token = NULL, password_reset_expiration = NULL WHERE id = ?",
		hashedPwd, userID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": i18n.T("password_update_error")})
		return
	}

	// Success response
	json.NewEncoder(w).Encode(map[string]string{"message": i18n.T("password_updated_successfully")})
}
