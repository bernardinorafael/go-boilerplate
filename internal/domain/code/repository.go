package code

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

func (r repo) Update(ctx context.Context, model model.Code) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		UPDATE codes
		SET
			active = :active,
			attempts = :attempts,
			used_at = :used_at,
			expires_at = :expires_at,
			updated_at = :updated_at
		WHERE id = :id
		AND user_id = :user_id
	`

	_, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return fault.New("failed to update code", fault.WithError(err))
	}

	return nil
}

func (r repo) GetByUserID(ctx context.Context, userID string) (*model.Code, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var code model.Code
	err := r.db.GetContext(
		ctx,
		&code,
		"SELECT * FROM codes WHERE user_id = $1 AND active = TRUE LIMIT 1",
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fault.New("code not found", fault.WithTag(fault.NotFound))
		}
		return nil, fault.New("failed to retrieve code", fault.WithError(err))
	}

	return &code, nil
}

func (r repo) InactivateAll(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		"UPDATE codes SET active = FALSE WHERE user_id = $1 AND active = TRUE",
		userID,
	)
	if err != nil {
		return fault.New("failed to update codes", fault.WithError(err))
	}

	return nil
}

func (r repo) Insert(ctx context.Context, model model.Code) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		INSERT INTO codes (
			id,
			user_id,
			code,
			active,
			attempts,
			used_at,
			expires_at,
			created_at,
			updated_at
		) VALUES (
			:id,
			:user_id,
			:code,
			:active,
			:attempts,
			:used_at,
			:expires_at,
			:created_at,
			:updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return fault.New("failed to insert product", fault.WithError(err))
	}

	return nil
}
