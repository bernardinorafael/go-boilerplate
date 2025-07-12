package dto

import (
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateCategory struct {
	Name string `json:"name"`
}

func (c CreateCategory) ValidateCreateCategory() error {
	return v.ValidateStruct(&c, v.Field(&c.Name, v.Required.Error("this field cannot be empty")))
}
