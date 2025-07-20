package dto

import (
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateProduct struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

func (p CreateProduct) Validate() error {
	return v.ValidateStruct(
		&p,
		v.Field(&p.Name, v.Required.Error("this field cannot be empty")),
		v.Field(&p.Price, v.Required.Error("this field cannot be empty or zero")),
	)
}

type UpdateProduct struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type ProductResponse struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Price   int64     `json:"price"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type CreateProductCategory struct {
	ProductID  string `json:"product_id"`
	CategoryID string `json:"category_id"`
}

func (p CreateProductCategory) Validate() error {
	return v.ValidateStruct(
		&p,
		v.Field(&p.CategoryID, v.Required.Error("this field cannot be empty")),
		v.Field(&p.ProductID, v.Required.Error("this field cannot be empty")),
	)
}
