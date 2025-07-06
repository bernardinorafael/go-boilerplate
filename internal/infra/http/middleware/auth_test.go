package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidateAccessToken(t *testing.T) {
	const (
		validToken = "mock_valid_token"
	)

	t.Run("should validate if Authorization is empty", func(t *testing.T) {
		got, err := validateAccessToken("")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "UNAUTHORIZED_ERROR: missing authentication token")
		assert.Empty(t, got)
	})

	t.Run("should validate if Authorization contains Bearer prefix", func(t *testing.T) {
		got, err := validateAccessToken(validToken)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "UNAUTHORIZED_ERROR: invalid auth token format")
		assert.Empty(t, got)
	})

	t.Run("should validate correct prefix", func(t *testing.T) {
		got, err := validateAccessToken("invalid " + validToken)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "UNAUTHORIZED_ERROR: invalid auth token format")
		assert.Empty(t, got)
	})

	t.Run("should validate correct prefix lowercase", func(t *testing.T) {
		got, err := validateAccessToken("bearer " + validToken)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "UNAUTHORIZED_ERROR: invalid auth token format")
		assert.Empty(t, got)
	})

	t.Run("should validate correct prefix uppercase", func(t *testing.T) {
		got, err := validateAccessToken("BEARER " + validToken)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "UNAUTHORIZED_ERROR: invalid auth token format")
		assert.Empty(t, got)
	})

	t.Run("should return token wihtout bearer prefix in success case", func(t *testing.T) {
		got, err := validateAccessToken("Bearer " + validToken)

		assert.Nil(t, err)
		assert.NotEmpty(t, got)
		assert.Equal(t, got, validToken)
	})
}
