package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypts a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hashedPassword), nil
}

// PasswordMatches compares a plain text password with an encrypted password
// Returns true if the password matches, false otherwise
func PasswordMatches(password, encrypted string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	return err == nil
}
