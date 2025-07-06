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
		r.Post("/login", h.login)
		r.Post("/code/{userId}", h.verifyCode)
	})

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(m.WithAuth)
		r.Get("/me", h.getSignedUser)
	})
}

func (h handler) verifyCode(w http.ResponseWriter, r *http.Request) {
	var body dto.VerifyCode
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = body.ValidateCode()
	if err != nil {
		fault.NewHTTPError(w, fault.NewValidation("failed to validate body", err))
		return
	}

	userID := chi.URLParam(r, "userId")
	res, err := h.service.Verify(r.Context(), userID, body.Code)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, res)
}

func (h handler) login(w http.ResponseWriter, r *http.Request) {
	var body dto.Login
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = body.ValidateLogin()
	if err != nil {
		fault.NewHTTPError(w, fault.NewValidation("failed to validate body", err))
		return
	}

	err = h.service.Login(r.Context(), body.Email)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		fault.NewHTTPError(w, fault.NewValidation("failed to validate body", err))
		return
	}

	err = h.service.Register(r.Context(), body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
