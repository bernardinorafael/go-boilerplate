package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/dbutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"

	"github.com/bernardinorafael/gogem/pkg/fault"
	"github.com/lib/pq"
)

type service struct {
	log      logging.Logger
	userRepo Repository
}

func NewService(log logging.Logger, userRepo Repository) Service {
	return &service{
		log:      log,
		userRepo: userRepo,
	}
}

func (s *service) CreateUser(ctx context.Context, input dto.CreateUser) error {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return fault.NewBadRequest("failed to get user by email")
	}
	if user != nil {
		return fault.NewConflict("e-mail already taken")
	}

	newUser, err := New(input.Name, input.Username, input.Email, input.Password)
	if err != nil {
		return fault.NewUnprocessableEntity("failed to create user entity")

	}
	model := newUser.ToModel()

	if err = s.userRepo.Insert(ctx, model); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505 is the code for unique constraint violation
			field := dbutil.ExtractFieldFromDetail(pqErr.Detail)
			return fault.NewConflict(fmt.Sprintf("%s already taken", field))
		}
		return fault.NewBadRequest("failed to insert user")
	}

	return nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	}
	if user == nil {
		return nil, fault.NewNotFound("user not found")
	}

	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	}
	if user == nil {
		return nil, fault.NewNotFound("user not found")
	}

	return user, nil
}
