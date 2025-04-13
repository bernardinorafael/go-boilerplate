package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/dbutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"
	"github.com/bernardinorafael/go-boilerplate/pkg/pagination"
	"github.com/lib/pq"
)

type ServiceConfig struct {
	ProductRepo Repository

	Metrics *metric.Metric
	Cache   *cache.Cache
}

type service struct {
	productRepo Repository
	metrics     *metric.Metric
	cache       *cache.Cache
}

func NewService(c ServiceConfig) Service {
	return &service{
		productRepo: c.ProductRepo,
		metrics:     c.Metrics,
		cache:       c.Cache,
	}
}

func (s service) GetProducts(ctx context.Context, search dto.SearchParams) (*pagination.Paginated[dto.ProductResponse], error) {
	records, totalItems, err := s.productRepo.GetAll(ctx, search)
	if err != nil {
		s.metrics.RecordError("products", "get-all-products")
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
	product, err := New(input.Name, input.Price)
	if err != nil {
		slog.Error("error create product entity", "error", err)
		return err // Error is already handled by the entity
	}

	err = s.productRepo.Insert(ctx, product.Model())
	if err != nil {
		var pqErr *pq.Error
		// 23505 is the code for unique contraint violation
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			field := dbutil.ExtractFieldFromDetail(pqErr.Detail)
			return fault.NewConflict(fmt.Sprintf("field: %s already taken", field))
		}

		s.metrics.RecordError("products", "insert-product")
		return fault.NewBadRequest("failed to insert product")
	}

	return nil
}

func (s service) GetProductByID(ctx context.Context, productID string) (*dto.ProductResponse, error) {
	var cachedProduct *dto.ProductResponse

	err := s.cache.GetStruct(ctx, fmt.Sprintf("prod:%s", productID), &cachedProduct)
	if err != nil {
		switch {
		case fault.GetTag(err) == fault.CACHE_MISS:
			s.metrics.RecordCacheMiss("product-service")
			slog.Info("cache:miss product not found")
		default:
			slog.Error("failed to query product from cache")
		}
	}

	if cachedProduct != nil {
		s.metrics.RecordCacheHit("product-service")
		return cachedProduct, nil
	}

	productRecord, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		s.metrics.RecordError("products", "get-by-id")
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

	cacheKey := fmt.Sprintf("prod:%s", product.ID)
	err = s.cache.SetStruct(ctx, cacheKey, product, time.Minute*30)
	if err != nil {
		slog.Error("failed to caching product")
	}

	return &product, nil
}

func (s service) DeleteProduct(ctx context.Context, productID string) error {
	productRecord, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		s.metrics.RecordError("products", "get-by-id")
		return fault.NewBadRequest("failed to retrieve product by ID")
	} else if productRecord == nil {
		return fault.NewNotFound("product not found")
	}

	err = s.productRepo.Delete(ctx, productID)
	if err != nil {
		s.metrics.RecordError("products", "delete-product")
		return fault.NewBadRequest("failed to delete product")
	}

	return nil
}

func (s service) UpdateProduct(ctx context.Context, productID string, input dto.UpdateProduct) error {
	productRecord, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		s.metrics.RecordError("products", "get-by-id")
		return fault.NewBadRequest("failed to retrieve product by ID")
	} else if productRecord == nil {
		return fault.NewNotFound("product not found")
	}

	p := NewFromModel(*productRecord)
	p.ChangeName(input.Name)
	p.ChangePrice(input.Price)

	err = s.productRepo.Update(ctx, p.Model())
	if err != nil {
		s.metrics.RecordError("products", "update-product")
		return fault.NewBadRequest("failed to update product")
	}

	return nil
}
