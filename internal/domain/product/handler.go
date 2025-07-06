package product

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

	r.Route("/api/v1/products", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/", h.createProduct)
		r.Get("/", h.getProducts)
		r.Get("/{productId}", h.getProduct)
		r.Delete("/{productId}", h.deleteProduct)
		r.Patch("/{productId}", h.updateProduct)
	})
}

func (h handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")

	var body dto.UpdateProduct
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = h.service.UpdateProduct(r.Context(), productID, body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")

	err := h.service.DeleteProduct(r.Context(), productID)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h handler) getProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")

	product, err := h.service.GetProductByID(r.Context(), productID)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, product)
}

func (h handler) getProducts(w http.ResponseWriter, r *http.Request) {
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

	products, err := h.service.GetProducts(r.Context(), search)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, products)
}

func (h handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateProduct
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = body.Validate()
	if err != nil {
		fault.NewHTTPError(w, fault.NewValidation("failed to validate request body", err))
		return
	}

	_, err = h.service.CreateProduct(r.Context(), body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
