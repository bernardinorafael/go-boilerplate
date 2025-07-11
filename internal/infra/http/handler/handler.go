package handler

import (
	"github.com/bernardinorafael/go-boilerplate/internal/domain/product"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/user"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/container"
	"github.com/go-chi/chi/v5"
)

func RegisterHandler(r *chi.Mux, c *container.Container) {
	product.NewHandler(c.ProductService, c.Config.JWTSecretKey).Register(r)
	user.NewHandler(c.UserService, c.Config.JWTSecretKey).Register(r)
}
