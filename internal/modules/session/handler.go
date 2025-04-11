package session

import (
	"net/http"
	"sync"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
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
		// Public
		r.Post("/refresh", h.handleRenewToken)
	})
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
