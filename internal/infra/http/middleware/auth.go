package middleware

import (
	"context"
	"gulg/pkg/fault"
	"gulg/pkg/token"
	"net/http"
	"strings"
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
			fault.NewHTTPResponse(w, fault.NewUnauthorized("access token not provided", nil))
			return
		}

		claims, err := token.Verify(m.secretKey, accessToken)
		if err != nil {
			if strings.Contains(err.Error(), "token has expired") {
				fault.NewHTTPResponse(w, fault.NewUnauthorized("token has expired", err))
				return
			}
			fault.NewHTTPResponse(w, fault.NewUnauthorized("invalid access token", err))
			return
		}

		ctx := context.WithValue(r.Context(), AuthKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
