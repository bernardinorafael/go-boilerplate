package session

import (
	"net/http"
	"sync"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/token"
	httputil "github.com/bernardinorafael/gogem/pkg/httputil"
	"github.com/go-chi/chi"
)

var (
	instance *handler
	once     sync.Once
)

type handler struct {
	sessionService Service
	secretKey      string
}

func NewHandler(sessionService Service, secretKey string) *handler {
	once.Do(func() {
		instance = &handler{
			sessionService: sessionService,
			secretKey:      secretKey,
		}
	})
	return instance
}

func (h handler) Register(r *chi.Mux) {
	m := middleware.NewWithAuth(h.secretKey)

	r.Route("/api/v1/sessions", func(r chi.Router) {
		// Private
		r.With(m.WithAuth).Get("/", h.handleGetSessions)
		r.With(m.WithAuth).Get("/me", h.handleGetSignedSession)
		// Public
		r.Post("/refresh", h.handleRenewToken)
	})
}

func (h handler) handleGetSignedSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		fault.NewHTTPError(w, fault.NewUnauthorized("invalid access token"))
		return
	}

	res, err := h.sessionService.GetSessionByUserID(ctx, c.UserID)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, res)
}

func (h handler) handleRenewToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	res, err := h.sessionService.RenewAccessToken(ctx, body.RefreshToken)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, res)
}

func (h handler) handleGetSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessions, err := h.sessionService.GetAllSessions(ctx)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, sessions)
}
