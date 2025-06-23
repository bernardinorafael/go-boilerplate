package user

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

func (r repo) GetByID(ctx context.Context, id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to get user by id", fault.WithError(err))
	}

	return &user, nil
}

func (r repo) Insert(ctx context.Context, model model.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var query = `
		INSERT INTO users (
			id,
			name,
			email,
			created,
			updated
		) VALUES (
			:id,
			:name,
			:email,
			:created,
			:updated
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return fault.New("failed to insert user", fault.WithError(err))
	}

	return nil
}
