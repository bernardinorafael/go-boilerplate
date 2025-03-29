package user

import (
	"context"
	"errors"
	"gulg/internal/_shared/dto"
	"gulg/internal/infra/database/model"
	"gulg/pkg/fault"
	"strings"

	"github.com/lib/pq"
)

type service struct {
	userRepo Repository
}

func NewService(userRepo Repository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) CreateUser(ctx context.Context, input dto.CreateUser) error {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return fault.NewBadRequest("failed to get user by email", err)
	}
	if user != nil {
		return fault.NewConflict("e-mail already taken", nil)
	}

	newUser, err := New(input.Name, input.Username, input.Email, input.Password)
	if err != nil {
		return fault.NewUnprocessableEntity("failed to create user entity", err)
	}
	model := newUser.ToModel()

	if err = s.userRepo.Insert(ctx, model); err != nil {
		var pqErr *pq.Error
		// 23505 is the code for unique constraint violation
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch {
			case strings.Contains(pqErr.Detail, "username"):
				return fault.NewConflict("username already taken", nil)
			case strings.Contains(pqErr.Detail, "email"):
				return fault.NewConflict("username already taken", nil)
			default:
				return fault.NewConflict("failed to insert user", nil)
			}
		}
		return fault.NewBadRequest("failed to insert user", err)
	}

	return nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user", err)
	}
	if user == nil {
		return nil, fault.NewNotFound("user not found", err)
	}

	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user", err)
	}
	if user == nil {
		return nil, fault.NewNotFound("user not found", err)
	}

	return user, nil
}
