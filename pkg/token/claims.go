package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	// TODO: Fill with system claims
	jwt.RegisteredClaims
}

func NewClaims(accId, userId, username string, duration time.Duration) (*Claims, error) {
	// claims := &Claims{
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		Subject:   username,
	// 		IssuedAt:  jwt.NewNumericDate(time.Now()),
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
	// 	},
	// }
	// return claims, nil
	return &Claims{}, nil
}

func (a *Claims) Valid() error {
	if time.Now().After(a.ExpiresAt.Time) {
		return errors.New("token has expired")
	}
	return nil
}
