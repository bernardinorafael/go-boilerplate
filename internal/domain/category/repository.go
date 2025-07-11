package category

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

func (r repo) Insert(ctx context.Context, category model.Category) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		INSERT INTO categories (
			id,
			name,
			slug,
			active,
			created_at,
			updated_at
		) VALUES (
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
		UPDATE categories
		SET
			name = :name,
			slug = :slug,
			active = :active,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		return fault.New("failed to update category", fault.WithError(err))
	}

	return nil
}

func (r repo) FindByID(ctx context.Context, categoryID string) (*model.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		SELECT * FROM categories
		WHERE id = :id
		AND deleted_at IS NULL
	`

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
