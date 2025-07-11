package category

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
)

type Repository interface {
	Insert(ctx context.Context, model model.Category) error
	Update(ctx context.Context, model model.Category) error
	FindByID(ctx context.Context, categoryID string) (*model.Category, error)
}
