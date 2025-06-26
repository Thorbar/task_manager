package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"task-manager/backend-go/db"
	"task-manager/backend-go/internal/i18n"
	"task-manager/backend-go/models"
	"unicode"
)

// isValidPassword checks password strength according to security rules:
// minimum 8 characters, at least one uppercase, one lowercase,
// one digit, and one special character.
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char), unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

/* RegisterHandler processes user registration HTTP requests.
 Validates method, parses input, checks password strength,
 registers the user, and returns appropriate responses.
*/
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("method_not_allowed"),
		})
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("invalid_request_data"),
		})
		return
	}

	if !isValidPassword(req.Password) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("password_not_secure"),
		})
		return
	}

	service := NewService(db.DB)
	if err := service.RegisterUser(context.Background(), &req); err != nil {
		log.Println("User registration error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("user_already_exists"),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": i18n.T("user_registered"),
	})
}
