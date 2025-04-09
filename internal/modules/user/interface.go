package user

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
)

type Repository interface {
	Insert(ctx context.Context, user model.User) error
	GetByID(ctx context.Context, userId string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Delete(ctx context.Context, userId string) error
}

type Service interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
	CreateUser(ctx context.Context, input dto.CreateUser) error
}
