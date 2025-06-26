package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	"task-manager/backend-go/internal/i18n"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware verifies JWT token and protects routes.
// Extracts user ID from token and adds it to the request context.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, i18n.T("error.token_not_provided"), http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, i18n.T("error.invalid_token_format"), http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		jwtKey := []byte(os.Getenv("JWT_SECRET"))

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, i18n.T("error.invalid_token"), http.StatusUnauthorized)
			return
		}

		// Extract user_id from claims as float64, then convert to int
		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			http.Error(w, i18n.T("error.invalid_token_claims"), http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		// Add userID to context for downstream handlers
		ctx := contextWithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// contextKey type is used to define context keys for type safety
type contextKey string

const userIDKey = contextKey("user_id")

// contextWithUserID adds the userID to the context
func contextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext retrieves the userID from the context
func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}
