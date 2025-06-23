package user

import (
	"net/http"
	"sync"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/httputil"
	"github.com/go-chi/chi"
)

var (
	instance *handler
	once     sync.Once
)

type handler struct {
	service   Service
	secretKey string
}

func NewHandler(service Service, secretKey string) *handler {
	once.Do(func() {
		instance = &handler{
			service:   service,
			secretKey: secretKey,
		}
	})
	return instance
}

func (h handler) Register(r *chi.Mux) {
	m := middleware.NewWithAuth(h.secretKey)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", h.registerUser)
	})

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Get("/me", h.getSignedUser)
	})
}

func (h handler) getSignedUser(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetSignedUser(r.Context())
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, res)
}

func (h handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateUser
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = body.Validate()
	if err != nil {
		fault.NewHTTPError(
			w, fault.New(
				"failed to validate request body",
				fault.WithTag(fault.ValidationError),
				fault.WithHTTPCode(http.StatusUnprocessableEntity),
				fault.WithValidationError(err),
			),
		)
		return
	}

	res, err := h.service.Register(r.Context(), body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, res)
}
