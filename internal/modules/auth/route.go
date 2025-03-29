package auth

import (
	"gulg/internal/_shared/dto"
	"gulg/pkg/fault"
	"gulg/pkg/httputil"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

var (
	instance *handler
	Once     sync.Once
)

type handler struct {
	authService Service
}

func NewHandler(authService Service) *handler {
	Once.Do(func() {
		instance = &handler{
			authService: authService,
		}
	})
	return instance
}

func (h handler) Register(r *chi.Mux) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)
	})
}

func (h handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body dto.CreateUser
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPResponse(w, err)
		return
	}

	err = h.authService.Register(ctx, body)
	if err != nil {
		fault.NewHTTPResponse(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusCreated)
}

func (h handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPResponse(w, err)
		return
	}

	res, err := h.authService.Login(ctx, body.Email, body.Password, r.RemoteAddr, r.UserAgent())
	if err != nil {
		fault.NewHTTPResponse(w, err)
		return
	}

	// handle login here

	httputil.WriteJSON(w, http.StatusOK, res)
}
