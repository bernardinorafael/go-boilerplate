package session

import (
	"github.com/go-chi/chi"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h handler) Register(r *chi.Mux) {
	r.Route("/v1/sessions", func(r chi.Router) {
		// routes here
	})
}
