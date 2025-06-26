package user

import (
	"context"
	"encoding/json"
	"net/http"
	"task-manager/backend-go/db"
	"task-manager/backend-go/internal/i18n"
	"task-manager/backend-go/models"
)

/* LoginHandler handles HTTP login requests.
It validates the method, decodes the request body,
calls the login service, and returns a JWT token on success.
*/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("method_not_allowed"),
		})
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("invalid_request_data"),
		})
		return
	}

	service := NewService(db.DB)

	token, err := service.LoginUser(context.Background(), &req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("login_failed"),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": i18n.T("login_success"),
		"token":   token,
	})
}
