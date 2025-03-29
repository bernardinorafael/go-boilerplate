package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewClaims(userId string, duration time.Duration) (*Claims, error) {
	return &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func (a *Claims) Valid() error {
	if time.Now().After(a.ExpiresAt.Time) {
		return errors.New("token has expired")
	}
	return nil
}
