package auth

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
)

type Service interface {
	Register(ctx context.Context, input dto.CreateUser) error
	Login(ctx context.Context, email, password, ip, agent string) (*dto.LoginResponse, error)
	GetSignedUser(ctx context.Context) (*dto.UserResponse, error)
	Activate(ctx context.Context, userId string) error
	Logout(ctx context.Context) error
}
