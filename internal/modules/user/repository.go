package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"

	"github.com/bernardinorafael/gogem/pkg/fault"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r repo) Update(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var query = `
		UPDATE users
		SET
			name = :name,
			username = :username,
			email = :email,
			password = :password,
			avatar_url = :avatar_url,
			enabled = :enabled,
			locked = :locked,
			updated = :updated
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fault.New("failed to update user", fault.WithError(err))
	}

	return nil
}

func (r repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user by email", fault.WithError(err))
	}

	return &user, nil
}

func (r repo) Delete(ctx context.Context, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return fault.New("failed to delete user", fault.WithError(err))
	}

	return nil
}

func (r repo) GetByID(ctx context.Context, userId string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1 LIMIT 1", userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user", fault.WithError(err))
	}

	return &user, nil
}

func (r repo) Insert(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var query = `
		INSERT INTO users (
			id,
			name,
			username,
			email,
			password,
			avatar_url,
			enabled,
			locked,
			created,
			updated
		) VALUES (
			:id,
			:name,
			:username,
			:email,
			:password,
			:avatar_url,
			:enabled,
			:locked,
			:created,
			:updated
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fault.New("failed to insert user", fault.WithError(err))
	}

	return nil
}
