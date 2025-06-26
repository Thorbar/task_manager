package main

import (
	"log"
	"net/http"

	"task-manager/backend-go/config"
	"task-manager/backend-go/db"
	"task-manager/backend-go/internal/auth"
	"task-manager/backend-go/internal/i18n"
	"task-manager/backend-go/internal/task"
	"task-manager/backend-go/internal/user"

	"github.com/rs/cors"
)

func main() {
	// Load configuration from .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("%s: %v", i18n.T("error.config.load"), err)
	}

	// Connect to MySQL database
	if err := db.ConnectDB(cfg); err != nil {
		log.Fatalf("%s: %v", i18n.T("error.database.connection"), err)
	}

	// Load i18n messages (default: Spanish)
	if err := i18n.LoadMessages("es"); err != nil {
		log.Fatalf("%s: %v", i18n.T("error.i18n.load"), err)
	}

	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(i18n.T("response.api.home")))
	})
	mux.HandleFunc("/register", user.RegisterHandler)
	mux.HandleFunc("/login", user.LoginHandler)
	mux.HandleFunc("/forgot-password", user.ForgotPasswordHandler)
	mux.HandleFunc("/reset-password", user.ResetPasswordHandler)

	// Protected task and user routes
	mux.HandleFunc("/api/tasks/", auth.AuthMiddleware(task.TasksRouter))
	mux.HandleFunc("/api/user/", auth.AuthMiddleware(user.UserRouter))

	// Apply CORS policy
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	log.Printf("%s http://localhost%s", i18n.T("server.start"), cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, handler))
}
