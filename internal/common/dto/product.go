package dto

import "time"

type CreateProduct struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
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
