package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewRepo(db *sqlx.DB, timeout time.Duration) *repo {
	return &repo{
		db:      db,
		timeout: timeout,
	}
}

func (r repo) InsertProductCategory(ctx context.Context, model model.ProductCategory) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		insert into product_categories (
			id,
			product_id,
			category_id,
			created_at
		) values (
			:id,
			:product_id,
			:category_id,
			:created_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return fault.New("failed to batch insert product categories", fault.WithError(err))
	}

	return nil
}

func (r repo) Delete(ctx context.Context, productID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", productID)
	if err != nil {
		return fault.New("failed to delete product", fault.WithError(err))
	}

	return nil
}

func (r repo) GetAll(ctx context.Context, search dto.SearchParams) ([]model.Product, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var products = make([]model.Product, 0)
	var skip = (search.Page - 1) * search.Limit

	var query = fmt.Sprintf(
		`select p.*
		from products p
		where (
			to_tsvector('simple', p.name)
			@@ websearch_to_tsquery('simple', $1)
			or p.name ilike '%%' || $1 || '%%'
		)
		order by p.created %s
		limit $2 offset $3`,
		search.Sort,
	)

	err := r.db.SelectContext(ctx, &products, query, search.Term, search.Limit, skip)
	if err != nil {
		return nil, -1, fault.New("failed to get products", fault.WithError(err))
	}

	var count int
	err = r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM products")
	if err != nil {
		return nil, -1, fault.New("failed to count products", fault.WithError(err))
	}

	return products, count, nil
}

func (r repo) GetByID(ctx context.Context, productID string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var product model.Product
	err := r.db.GetContext(ctx, &product, "SELECT * FROM products WHERE id = $1", productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fault.New("product not found", fault.WithTag(fault.NotFound))
		}
		return nil, fault.New("failed to get product by id", fault.WithError(err))
	}

	return &product, nil
}

func (r repo) GetByName(ctx context.Context, name string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var product model.Product
	err := r.db.GetContext(ctx, &product, "SELECT * FROM products WHERE name = $1", name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fault.New("product not found", fault.WithTag(fault.NotFound))
		}
		return nil, fault.New("failed to get product by name", fault.WithError(err))
	}

	return &product, nil
}

func (r repo) Insert(ctx context.Context, product model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		insert into products (
			id,
			name,
			price,
			created,
			updated
		) values (
			:id,
			:name,
			:price,
			:created,
			:updated
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		return fault.New("failed to insert product", fault.WithError(err))
	}

	return nil
}

func (r repo) Update(ctx context.Context, product model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		update products
		set
			name = :name,
			price = :price,
			updated = :updated
		where id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		return fault.New("failed to update product", fault.WithError(err))
	}

	return nil
}
