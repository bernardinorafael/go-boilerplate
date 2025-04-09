package auth

import (
	"net/http"
	"sync"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/token"

	httputil "github.com/bernardinorafael/gogem/pkg/httputil"
	"github.com/go-chi/chi"
)

var (
	instance *handler
	Once     sync.Once
)

type handler struct {
	authService Service
	secretKey   string
}

func NewHandler(authService Service, secretKey string) *handler {
	Once.Do(func() {
		instance = &handler{
			authService: authService,
			secretKey:   secretKey,
		}
	})
	return instance
}

func (h handler) Register(r *chi.Mux) {
	m := middleware.NewWithAuth(h.secretKey)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.With(m.WithAuth).Get("/me", h.handleGetSigned)
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)
	})
}

func (h handler) handleGetSigned(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		fault.NewHTTPError(w, fault.NewUnauthorized("access token not found"))
		return
	}

	user, err := h.authService.GetSigned(ctx, userId.ID)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, user)
}

func (h handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body dto.CreateUser
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = h.authService.Register(ctx, body)
	if err != nil {
		fault.NewHTTPError(w, err)
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
		fault.NewHTTPError(w, err)
		return
	}

	res, err := h.authService.Login(ctx, body.Email, body.Password, r.RemoteAddr, r.UserAgent())
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, res)
}
