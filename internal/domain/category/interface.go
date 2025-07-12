package category

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/pagination"
)

type Repository interface {
	Insert(ctx context.Context, model model.Category) error
	Update(ctx context.Context, model model.Category) error
	FindByID(ctx context.Context, categoryID string) (*model.Category, error)
	FindByName(ctx context.Context, name string) (*model.Category, error)
	FindAll(ctx context.Context, search dto.SearchParams) ([]model.Category, int, error)
}

type Service interface {
	Create(ctx context.Context, input dto.CreateCategory) error
	Delete(ctx context.Context, categoryID string) error
	GetByID(ctx context.Context, categoryID string) (*dto.CategoryResponse, error)
	FindAll(ctx context.Context, search dto.SearchParams) (*pagination.Paginated[dto.CategoryResponse], error)
}
