package dto

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type SearchParams struct {
	Term  string `json:"term,omitempty"`
	Sort  string `json:"sort,omitempty"`
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

func (s SearchParams) Validate() error {
	return v.ValidateStruct(
		&s,
		v.Field(
			&s.Sort,
			v.Required.Error("this param cannot be empty"),
			v.In("DESC", "ASC").Error("this field must be one of the following values - ASC or DESC"),
		),
		v.Field(
			&s.Page,
			v.Required.Error("this param cannot be empty or zero"),
			v.Min(1).Error("this field must be greater than 0"),
		),
		v.Field(
			&s.Limit,
			v.Required.Error("this param cannot be empty or zero"),
			v.In(10, 20, 50, 100).Error("this field must be one of the following values - 10, 20, 50, or 100"),
		),
	)
}
