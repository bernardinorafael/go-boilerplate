package session

import (
	"github.com/go-chi/chi"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h handler) Register(r *chi.Mux) {
	r.Route("/api/v1/sessions", func(r chi.Router) {
		// handle here
	})
}
