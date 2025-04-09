package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r repo) GetActiveByUserID(ctx context.Context, userId string) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var session model.Session
	err := r.db.GetContext(
		ctx,
		&session,
		"SELECT * FROM sessions WHERE user_id = $1 AND active = true LIMIT 1",
		userId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	return &session, nil
}

func (r repo) Delete(ctx context.Context, sessionId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", sessionId)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r repo) GetByID(ctx context.Context, sessionId string) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var session model.Session
	err := r.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE id = $1 limit 1", sessionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	return &session, nil
}

func (r repo) Insert(ctx context.Context, session model.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var query = `
		INSERT INTO sessions (
			id,
			user_id,
			agent,
			ip_address,
			refresh_token,
			active,
			expires,
			created,
			updated
		)	VALUES (
			:id,
			:user_id,
			:agent,
			:ip_address,
			:refresh_token,
			:active,
			:expires,
			:created,
			:updated
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	return nil
}
