package config

import (
	"fmt"
	"os"
	"task-manager/backend-go/internal/i18n"

	"github.com/joho/godotenv"
)

// Config holds application environment variables
type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	Port       string
	JWTSecret  string
}

// Load reads environment variables and returns a Config instance
func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Valor por defecto
	}

	// Carga .env (base) + .env.[entorno] (específico)
	err := godotenv.Load(".env", ".env."+env)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.T("error.i18n.load"), err)
	}

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		Port:       getPort(),
	}, nil
}

// getPort returns the server port, defaulting to :8080 if not set
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}
	return ":" + port
}
