package product

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/dbutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"
	"github.com/bernardinorafael/go-boilerplate/pkg/pagination"
)

type ServiceConfig struct {
	ProductRepo Repository

	Log     *log.Logger
	Metrics *metric.Metric
	Cache   *cache.Cache
}

type service struct {
	log     *log.Logger
	repo    Repository
	metrics *metric.Metric
	cache   *cache.Cache
}

func NewService(c ServiceConfig) *service {
	return &service{
		log:     c.Log,
		repo:    c.ProductRepo,
		metrics: c.Metrics,
		cache:   c.Cache,
	}
}

func (s service) GetProducts(ctx context.Context, search dto.SearchParams) (*pagination.Paginated[dto.ProductResponse], error) {
	s.log.Debug(
		"trying to retrieve products with",
		"details", strings.Join(
			[]string{
				fmt.Sprintf("term: %s", search.Term),
				fmt.Sprintf("sort: %s", search.Sort),
				fmt.Sprintf("limit: %d", search.Limit),
				fmt.Sprintf("page: %d", search.Page),
			},
			"\n",
		),
	)

	records, totalItems, err := s.repo.GetAll(ctx, search)
	if err != nil {
		s.log.Error("failed to retrieve products", "err", err)
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

	s.log.Debug("products retrieved successfully", "totalItems", totalItems)

	paginatedResponse := pagination.New(products, totalItems, search.Page, search.Limit)
	return &paginatedResponse, nil
}

func (s service) CreateProduct(ctx context.Context, input dto.CreateProduct) (*dto.ProductResponse, error) {
	s.log.Debug(
		"trying to create product with",
		"name", input.Name,
		"price", input.Price,
	)

	p, err := NewEntity(input.Name, input.Price)
	if err != nil {
		s.log.Error("failed to create product", "err", err)
		return nil, err // Error is already handled by the entity
	}

	err = s.repo.Insert(ctx, p.Model())
	if err != nil {
		if err = dbutil.VerifyDuplicatedConstraintKey(err); err != nil {
			s.log.Error("duplicated product", "name", input.Name, "err", err)
			s.metrics.RecordError("user", "product-user")
			return nil, err // Error is already handled by the helper
		}
		s.metrics.RecordError("products", "insert-product")
		return nil, fault.NewBadRequest("failed to insert product")
	}

	s.log.Debug(
		"product created successfully",
		"details", strings.Join(
			[]string{
				fmt.Sprintf("id: %s", p.id),
				fmt.Sprintf("name: %s", p.name),
				fmt.Sprintf("price: %d", p.price),
			},
			"\n",
		),
	)

	res := dto.ProductResponse{
		ID:      p.id,
		Name:    p.name,
		Price:   p.price,
		Created: p.created,
		Updated: p.updated,
	}

	return &res, nil
}

func (s service) GetProductByID(ctx context.Context, productID string) (*dto.ProductResponse, error) {
	s.log.Debug("trying to retrieve product", "id", productID)

	var cachedProduct *dto.ProductResponse
	err := s.cache.GetStruct(ctx, fmt.Sprintf("prod:%s", productID), &cachedProduct)
	if err != nil {
		switch {
		case fault.GetTag(err) == fault.CacheMiss:
			s.log.Debug("cache miss for product", "id", productID)
			s.metrics.RecordCacheMiss("product")
		default:
			s.log.Error("failed to query product from cache", "err", err)
		}
	}

	if cachedProduct != nil {
		s.log.Debug("cache hit for product", "id", productID)
		s.metrics.RecordCacheHit("product")
		return cachedProduct, nil
	}

	record, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			return nil, fault.NewNotFound("product not found")
		default:
			s.metrics.RecordError("products", "get-by-id")
			return nil, fault.NewBadRequest("failed to retrieve product by ID")
		}
	}

	product := dto.ProductResponse{
		ID:      record.ID,
		Name:    record.Name,
		Price:   record.Price,
		Created: record.Created,
		Updated: record.Updated,
	}

	cacheKey := fmt.Sprintf("prod:%s", product.ID)
	err = s.cache.SetStruct(ctx, cacheKey, product, time.Minute*30)
	if err != nil {
		s.log.Error("failed to caching product", "err", err)
		s.metrics.RecordError("products", "cache-product")
	}
	s.log.Debug("product stored in cache", "cacheKey", cacheKey)

	return &product, nil
}

func (s service) DeleteProduct(ctx context.Context, productID string) error {
	s.log.Debug("trying to delete product", "id", productID)

	_, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			return fault.NewNotFound("product not found")
		default:
			s.metrics.RecordError("products", "get-by-id")
			return fault.NewBadRequest("failed to retrieve product by ID")
		}
	}

	err = s.repo.Delete(ctx, productID)
	if err != nil {
		s.log.Error("failed to delete product", "err", err)
		s.metrics.RecordError("products", "delete-product")
		return fault.NewBadRequest("failed to delete product")
	}

	cacheKey := fmt.Sprintf("prod:%s", productID)

	existInCache, err := s.cache.Has(ctx, cacheKey)
	if err != nil {
		s.log.Error("failed to search for product in cache", "err", err)
		return fault.NewInternalServerError("failed to search for product")
	}
	if existInCache {
		err := s.cache.Delete(ctx, cacheKey)
		if err != nil {
			s.log.Error("failed to delete product from cache", "err", err)
			return fault.NewInternalServerError("failed to delete product from cache")
		}

		s.log.Debug("product deleted from cache", "cacheKey", cacheKey)
	}

	s.log.Debug("product deleted successfully", "id", productID)

	return nil
}

func (s service) UpdateProduct(ctx context.Context, productID string, input dto.UpdateProduct) error {
	record, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		switch fault.GetTag(err) {
		case fault.NotFound:
			return fault.NewNotFound("product not found")
		default:
			s.metrics.RecordError("products", "get-by-id")
			return fault.NewBadRequest("failed to retrieve product by ID")
		}
	}

	entity := NewFromModel(*record)
	entity.ChangeName(input.Name)
	entity.ChangePrice(input.Price)

	err = s.repo.Update(ctx, entity.Model())
	if err != nil {
		s.metrics.RecordError("products", "update-product")
		return fault.NewBadRequest("failed to update product")
	}

	return nil
}
