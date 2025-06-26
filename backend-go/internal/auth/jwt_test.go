package auth_test

import (
	"testing"

	"task-manager/backend-go/internal/auth"
	"task-manager/backend-go/internal/i18n"

	"github.com/stretchr/testify/assert"
)

// TestGenerateJWT verifies JWT token creation with a valid user ID
func TestGenerateJWT(t *testing.T) {
	token, err := auth.GenerateJWT(123)
	assert.NoError(t, err, i18n.T("test.jwt.generate.no_error"))
	assert.NotEmpty(t, token, i18n.T("test.jwt.generate.not_empty"))
}

// TestParseToken verifies parsing of a valid JWT token and matching user ID
func TestParseToken(t *testing.T) {
	userID := 123

	// Generate token for testing parsing
	token, err := auth.GenerateJWT(userID)
	assert.NoError(t, err, i18n.T("test.jwt.parse.generate_no_error"))

	parsedUserID, err := auth.ParseToken(token)
	assert.NoError(t, err, i18n.T("test.jwt.parse.parse_no_error"))
	assert.Equal(t, userID, parsedUserID, i18n.T("test.jwt.parse.match_id"))
}
