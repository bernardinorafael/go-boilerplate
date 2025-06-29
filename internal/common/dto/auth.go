package dto

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type AuthResponse struct {
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserResponse struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type Login struct {
	Email string `json:"email"`
}

func (l Login) ValidateLogin() error {
	return v.ValidateStruct(
		&l,
		v.Field(
			&l.Email,
			v.Required.Error("is a required field"),
			is.Email.Error("invalid email format"),
		),
	)
}

type VerifyCode struct {
	Code string `json:"code"`
}

func (c VerifyCode) ValidateCode() error {
	return v.ValidateStruct(
		&c,
		v.Field(
			&c.Code,
			v.Required.Error("is a required field"),
			v.Length(6, 6).Error("code must be exactly 6 characters long"),
		),
	)
}
