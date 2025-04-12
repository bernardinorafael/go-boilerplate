package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/pkg/dbutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"
	"github.com/bernardinorafael/go-boilerplate/pkg/pagination"
	"github.com/lib/pq"
)

type service struct {
	log         logging.Logger
	productRepo Repository
}

func NewService(log logging.Logger, productRepo Repository) Service {
	return &service{
		log:         log,
		productRepo: productRepo,
	}
}

func (s service) GetProducts(ctx context.Context, search dto.SearchParams) (*pagination.Paginated[dto.ProductResponse], error) {
	records, totalItems, err := s.productRepo.GetAll(ctx, search)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve products")
	}

	var products = make([]dto.ProductResponse, len(records))
	for i, p := range records {
		products[i] = dto.ProductResponse{
			ID:      p.ID,
			Name:    p.Name,
			Price:   p.Price,
			Created: p.Created,
			Updated: p.Updated,
		}
	}

	res := pagination.New(products, totalItems, search.Page, search.Limit)
	return &res, nil
}

func (s service) CreateProduct(ctx context.Context, input dto.CreateProduct) error {
	product := New(input.Name, input.Price)

	err := s.productRepo.Insert(ctx, product.Model())
	if err != nil {
		var pqErr *pq.Error
		// 23505 is the code for unique contraint violation
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			field := dbutil.ExtractFieldFromDetail(pqErr.Detail)
			return fault.NewConflict(fmt.Sprintf("field: %s already taken", field))
		}
		return fault.NewBadRequest("failed to insert product")
	}

	return nil
}

func (s service) GetProductByID(ctx context.Context, productID string) (*dto.ProductResponse, error) {
	productRecord, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve product by ID")
	} else if productRecord == nil {
		return nil, fault.NewNotFound("product not found")
	}

	product := dto.ProductResponse{
		ID:      productRecord.ID,
		Name:    productRecord.Name,
		Price:   productRecord.Price,
		Created: productRecord.Created,
		Updated: productRecord.Updated,
	}

	return &product, nil
}

func (s service) DeleteProduct(ctx context.Context, productID string) error {
	productRecord, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fault.NewBadRequest("failed to retrieve product by ID")
	} else if productRecord == nil {
		return fault.NewNotFound("product not found")
	}

	err = s.productRepo.Delete(ctx, productID)
	if err != nil {
		return fault.NewBadRequest("failed to delete product")
	}

	return nil
}

func (s service) UpdateProduct(ctx context.Context, productID string, input dto.UpdateProduct) error {
	productRecord, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fault.NewBadRequest("failed to retrieve product by ID")
	} else if productRecord == nil {
		return fault.NewNotFound("product not found")
	}

	product := NewFromModel(*productRecord)

	product.ChangeName(input.Name)
	product.ChangePrice(input.Price)

	err = s.productRepo.Update(ctx, product.Model())
	if err != nil {
		return fault.NewBadRequest("failed to update product")
	}

	return nil
}
