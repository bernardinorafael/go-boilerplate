package user

import (
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/crypto"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
)

type user struct {
	id         string
	name       string
	username   string
	email      string
	password   string
	avatar_url *string
	enabled    bool
	locked     bool
	created    time.Time
	updated    time.Time
}

func New(name, username, email, pass string) (*user, error) {
	hashedPass, err := crypto.HashPassword(pass)
	if err != nil {
		return nil, fault.New("failed to hash password", fault.WithError(err))
	}

	u := user{
		id:         uid.New("user"),
		name:       name,
		username:   username,
		email:      email,
		password:   hashedPass,
		avatar_url: nil,
		enabled:    false,
		locked:     false,
		created:    time.Now(),
		updated:    time.Now(),
	}

	if err := u.validate(); err != nil {
		return nil, fault.New(
			"failed to create user entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &u, nil
}

func (u *user) ToModel() model.User {
	return model.User{
		ID:        u.id,
		Name:      u.name,
		Username:  u.username,
		Email:     u.email,
		Password:  u.password,
		AvatarURL: u.avatar_url,
		Enabled:   u.enabled,
		Locked:    u.locked,
		Created:   u.created,
		Updated:   u.updated,
	}
}

func (u *user) validate() error {
	if u.name == "" {
		return fault.New("user name is required")
	}
	if u.password == "" {
		return fault.New("password is required")
	}
	if u.email == "" {
		return fault.New("email is required")
	}
	if u.username == "" {
		return fault.New("username is required")
	}

	return nil
}
