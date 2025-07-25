package category

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

func (r repo) FindAll(ctx context.Context, search dto.SearchParams) ([]model.Category, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var categories = make([]model.Category, 0)
	var skip = (search.Page - 1) * search.Limit

	var query = fmt.Sprintf(
		`select c.*
		from categories c
		where (
			to_tsvector('simple', c.name)
			@@ websearch_to_tsquery('simple', $1)
			or c.name ilike '%%' || $1 || '%%'
		)
		and c.deleted_at is null
		order by c.created_at %s
		limit $2 offset $3`,
		search.Sort,
	)

	err := r.db.SelectContext(ctx, &categories, query, search.Term, search.Limit, skip)
	if err != nil {
		return nil, -1, fault.New("failed to get categories", fault.WithError(err))
	}

	var count int
	err = r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM categories")
	if err != nil {
		return nil, -1, fault.New("failed to count categories", fault.WithError(err))
	}

	return categories, count, nil
}

func (r repo) Insert(ctx context.Context, category model.Category) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		insert into categories (
			id,
			name,
			slug,
			active,
			created_at,
			updated_at
		) values (
			:id,
			:name,
			:slug,
			:active,
			:created_at,
			:updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		return fault.New("failed to insert category", fault.WithError(err))
	}

	return nil
}

func (r repo) Update(ctx context.Context, category model.Category) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		update categories
		set
			name = :name,
			slug = :slug,
			active = :active,
			updated_at = :updated_at,
			deleted_at = :deleted_at
		where id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		return fault.New("failed to update category", fault.WithError(err))
	}

	return nil
}

func (r repo) FindByName(ctx context.Context, name string) (*model.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `select * from categories where name = $1`

	var category model.Category
	err := r.db.GetContext(ctx, &category, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fault.New("category not found", fault.WithTag(fault.NotFound))
		}
		return nil, fault.New("failed to get category by id", fault.WithError(err))
	}

	return &category, nil
}

func (r repo) FindByID(ctx context.Context, categoryID string) (*model.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `select * from categories where id = $1`

	var category model.Category
	err := r.db.GetContext(ctx, &category, query, categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fault.New("category not found", fault.WithTag(fault.NotFound))
		}
		return nil, fault.New("failed to get category by id", fault.WithError(err))
	}

	return &category, nil
}
