package user

import (
	"net/http"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/model"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	minNameLength = 3
	maxNameLength = 100
)

type user struct {
	id      string
	name    string
	email   string
	created time.Time
	updated time.Time
}

func NewEntity(name, email string) (*user, error) {
	u := user{
		id:      uid.New("user"),
		name:    name,
		email:   email,
		created: time.Now(),
		updated: time.Now(),
	}

	if err := u.validate(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (u user) validate() error {
	err := v.ValidateStruct(
		&u,
		v.Field(
			&u.name,
			v.Required.Error("this field is required"),
			v.Length(minNameLength, maxNameLength).Error("name must be between 3 and 100 characters"),
		),
		v.Field(
			&u.email,
			v.Required.Error("this field is required"),
			is.Email.Error("invalid email format"),
		),
	)
	if err != nil {
		return fault.New(
			"failed to validate product",
			fault.WithHTTPCode(http.StatusUnprocessableEntity),
			fault.WithTag(fault.ValidationError),
			fault.WithValidationError(err),
		)
	}

	return nil
}

func (u *user) Model() model.User {
	return model.User{
		ID:      u.id,
		Name:    u.name,
		Email:   u.email,
		Created: u.created,
		Updated: u.updated,
	}
}
