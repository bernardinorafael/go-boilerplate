package auth

import (
	"net/http"
	"sync"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"

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
		// Private
		r.With(m.WithAuth).Get("/me", h.handleGetSigned)
		r.With(m.WithAuth).Patch("/logout", h.handleLogout)
		// Public
		r.Get("/activate/{userId}", h.handleActivate)
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)
	})
}

func (h handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.authService.Logout(ctx)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusOK)
}

func (h handler) handleActivate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := chi.URLParam(r, "userId")

	err := h.authService.Activate(ctx, userId)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusOK)
}

func (h handler) handleGetSigned(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.authService.GetSignedUser(ctx)
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

	httputil.WriteSuccess(w, http.StatusAccepted)
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
