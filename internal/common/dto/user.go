package dto

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u CreateUser) Validate() error {
	return v.ValidateStruct(
		&u,
		v.Field(&u.Name, v.Required.Error("this field cannot be empty")),
		v.Field(&u.Email, v.Required.Error("this field cannot be empty")),
	)
}
