package category

import (
	"net/http"
	"sync"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/httputil"
	"github.com/go-chi/chi/v5"
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

	r.Route("/api/v1/categories", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/", h.createCategory)
		r.Get("/", h.getCategories)
		r.Delete("/{id}", h.deleteCategory)
		r.Get("/{id}", h.getCategory)
	})
}

func (h handler) getCategories(w http.ResponseWriter, r *http.Request) {
	var q = r.URL.Query()
	var search dto.SearchParams

	search.Term = httputil.ReadQueryString(q, "term")
	search.Sort = httputil.ReadQueryString(q, "sort")
	search.Limit = httputil.ReadQueryInt(q, "limit")
	search.Page = httputil.ReadQueryInt(q, "page")

	err := search.Validate()
	if err != nil {
		fault.NewHTTPError(w, fault.NewValidation("failed to validate query params", err))
		return
	}

	payload, err := h.service.FindAll(r.Context(), search)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, payload)
}

func (h handler) createCategory(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateCategory
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = body.ValidateCreateCategory()
	if err != nil {
		fault.NewHTTPError(w, fault.NewValidation("failed to validate request body", err))
		return
	}

	err = h.service.Create(r.Context(), body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), categoryID)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h handler) getCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")

	category, err := h.service.GetByID(r.Context(), categoryID)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, category)
}
