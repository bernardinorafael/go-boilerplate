package auth

import (
	"context"
	"gulg/internal/_shared/dto"
)

type Service interface {
	Register(ctx context.Context, input dto.CreateUser) error
	Login(ctx context.Context, email, password, ip, agent string) (*dto.LoginResponse, error)
}
