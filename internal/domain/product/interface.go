package product

import (
	"context"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/pagination"
)

type Repository interface {
	Insert(ctx context.Context, product model.Product) error
	GetAll(ctx context.Context, search dto.SearchParams) ([]model.Product, int, error)
	GetByID(ctx context.Context, productID string) (*model.Product, error)
	GetByName(ctx context.Context, name string) (*model.Product, error)
	Update(ctx context.Context, product model.Product) error
	Delete(ctx context.Context, productID string) error
	InsertProductCategory(ctx context.Context, model model.ProductCategory) error
}

type Service interface {
	CreateProduct(ctx context.Context, input dto.CreateProduct) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, productID string, input dto.UpdateProduct) error
	GetProducts(ctx context.Context, search dto.SearchParams) (*pagination.Paginated[dto.ProductResponse], error)
	GetProductByID(ctx context.Context, productID string) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, productID string) error
	AddProductCategory(ctx context.Context, productID string, categoryID string) error
}
