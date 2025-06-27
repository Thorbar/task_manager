package user

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"task-manager/backend-go/db"
	"task-manager/backend-go/internal/i18n"
	"task-manager/backend-go/models"
	"time"

	"github.com/joho/godotenv"
)

// Email configuration variables loaded from env files.
var (
	from     string
	password string
	url      string
	port     string
)

func init() {
	// Load environment variables depending on the app environment
	/*
	env := os.Getenv("APP_ENV") // "development" or "production"

	envFile := ".env.development"
	if env == "production" {
		envFile = ".env.production"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Could not load env file %s: %v", envFile, err)
	}
*/
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Could not load env file %s: ", err)
	}

	from = os.Getenv("GOOGLE_EMAIL")
	password = os.Getenv("GOOGLE_PWD")
	url = os.Getenv("PUBLIC_API_URL")
	port = os.Getenv("PUBLIC_API_PORT")

	if from == "" || password == "" || url == "" {
		log.Printf("Missing mandatory environment variables: GOOGLE_EMAIL, GOOGLE_PWD, PUBLIC_API_URL")
	}
}

// generateResetToken creates a secure random token for password reset.
func generateResetToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// sendResetEmail constructs and sends the password reset email with token link.
func sendResetEmail(email, token string) error {
	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	var resetURL string
	if port != "" {
		resetURL = fmt.Sprintf("%s:%s/resetPasswordRequest?token=%s", url, port, token)
	} else {
		resetURL = fmt.Sprintf("%s/resetPasswordRequest?token=%s", url, token)
	}

	// Set token expiration to 15 minutes from now and store in DB
	expiration := time.Now().Add(15 * time.Minute)
	_, err := db.DB.Exec("UPDATE users SET password_reset_token = ?, password_reset_expiration = ? WHERE email = ?", token, expiration, email)
	if err != nil {
		log.Printf("Error saving token in database: %v", err)
		return err
	}

	// Email content with i18n translations
	subject := i18n.T("forgot_email_subject")
	body := i18n.T("forgot_email_intro") + "\r\n\r\n" +
		i18n.T("forgot_email_instruction") + "\r\n" + resetURL + "\r\n\r\n" +
		i18n.T("forgot_email_ignore") + "\r\n\r\n" +
		i18n.T("forgot_email_team")

	message := []byte(subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	return nil
}

// ForgotPasswordHandler processes password reset requests via POST.
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("error_method_not_allowed"),
		})
		return
	}

	var user models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("error_invalid_data"),
		})
		return
	}

	err := db.DB.QueryRow("SELECT username, email FROM users WHERE email = ?", user.Email).Scan(&user.Username, &user.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("error_email_not_registered"),
		})
		return
	}

	token, err := generateResetToken()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("error_token_generation"),
		})
		return
	}


	if err := sendResetEmail(user.Email, token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": i18n.T("error_email_send_fail"),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": i18n.T("message_email_sent"),
	})
}
