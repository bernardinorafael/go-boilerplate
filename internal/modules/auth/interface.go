package auth

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
)

type Service interface {
	Register(ctx context.Context, input dto.CreateUser) error
	Login(ctx context.Context, email, password, ip, agent string) (*dto.LoginResponse, error)
	GetSigned(ctx context.Context, userId string) (*model.User, error)
}
