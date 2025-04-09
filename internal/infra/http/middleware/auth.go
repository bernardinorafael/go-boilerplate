package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/token"
)

type AuthKey struct{}

type middleware struct {
	secretKey string
}

func NewWithAuth(secretKey string) *middleware {
	return &middleware{
		secretKey: secretKey,
	}
}

func (m *middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")

		if len(accessToken) == 0 {
			fault.NewHTTPError(w, fault.NewUnauthorized("access token not provided"))
			return
		}

		claims, err := token.Verify(m.secretKey, accessToken)
		if err != nil {
			if strings.Contains(err.Error(), "token has expired") {
				fault.NewHTTPError(w, fault.NewUnauthorized("token has expired"))
				return
			}
			fault.NewHTTPError(w, fault.NewUnauthorized("invalid access token"))
			return
		}

		ctx := context.WithValue(r.Context(), AuthKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
