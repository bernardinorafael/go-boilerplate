package code

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
)

type Service interface {
	CreateCode(ctx context.Context, userID string) error
	VerifyCode(ctx context.Context, userID, code string) (bool, error)
}

type Repository interface {
	Insert(ctx context.Context, model model.Code) error
	InactivateAll(ctx context.Context, userID string) error
	GetByUserID(ctx context.Context, userID string) (*model.Code, error)
	Update(ctx context.Context, model model.Code) error
}
