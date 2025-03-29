package session

import (
	"context"
	"gulg/internal/infra/database/model"
)

type Repository interface {
	Insert(ctx context.Context, session model.Session) error
	GetByID(ctx context.Context, sessionId string) (*model.Session, error)
	Delete(ctx context.Context, sessionId string) error
}

type Service interface {
	CreateSession(ctx context.Context, userId, ip, agent, refreshToken string) (sessionId string, err error)
}
