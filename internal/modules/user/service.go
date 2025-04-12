package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/pkg/dbutil"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"

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

func (s service) CreateUser(ctx context.Context, input dto.CreateUser) (*dto.UserResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to get user by email")
	} else if userRecord != nil {
		return nil, fault.NewConflict("e-mail already taken")
	}

	newUser, err := New(input.Name, input.Username, input.Email, input.Password)
	if err != nil {
		return nil, fault.NewUnprocessableEntity("failed to create user entity")

	}
	model := newUser.Model()

	if err = s.userRepo.Insert(ctx, model); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505 is the code for unique constraint violation
			field := dbutil.ExtractFieldFromDetail(pqErr.Detail)
			return nil, fault.NewConflict(fmt.Sprintf("%s already taken", field))
		}
		return nil, fault.NewBadRequest("failed to insert user")
	}

	user := dto.UserResponse{
		ID:        model.ID,
		Name:      model.Name,
		Username:  model.Username,
		Email:     model.Email,
		AvatarURL: model.AvatarURL,
		Locked:    model.Locked,
		Created:   model.Created,
		Updated:   model.Updated,
	}

	return &user, nil
}

func (s service) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	}
	if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user := dto.UserResponse{
		ID:        userRecord.ID,
		Name:      userRecord.Name,
		Username:  userRecord.Username,
		Email:     userRecord.Email,
		AvatarURL: userRecord.AvatarURL,
		Locked:    userRecord.Locked,
		Created:   userRecord.Created,
		Updated:   userRecord.Updated,
	}

	return &user, nil
}

func (s service) GetUserByID(ctx context.Context, userId string) (*dto.UserResponse, error) {
	userRecord, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	} else if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user := dto.UserResponse{
		ID:        userRecord.ID,
		Name:      userRecord.Name,
		Username:  userRecord.Username,
		Email:     userRecord.Email,
		AvatarURL: userRecord.AvatarURL,
		Locked:    userRecord.Locked,
		Created:   userRecord.Created,
		Updated:   userRecord.Updated,
	}

	return &user, nil
}
