package user

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
)

type Repository interface {
	Insert(ctx context.Context, model model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type Service interface {
	Register(ctx context.Context, input dto.CreateUser) error
	GetSignedUser(ctx context.Context) (*dto.UserResponse, error)
	Login(ctx context.Context, email string) error
	Verify(ctx context.Context, userID, code string) (*dto.AuthResponse, error)
}
