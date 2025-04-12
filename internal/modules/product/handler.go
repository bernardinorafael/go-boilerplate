package product

import (
	"net/http"
	"sync"

	"github.com/bernardinorafael/go-boilerplate/internal/common/dto"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/go-chi/chi"

	httputil "github.com/bernardinorafael/gogem/pkg/httputil"
)

var (
	instance *handler
	once     sync.Once
)

type handler struct {
	productService Service
	secretKey      string
}

func NewHandler(productService Service, secretKey string) *handler {
	once.Do(func() {
		instance = &handler{
			productService: productService,
			secretKey:      secretKey,
		}
	})
	return instance
}

func (h handler) Register(r *chi.Mux) {
	m := middleware.NewWithAuth(h.secretKey)

	r.Route("/api/v1/products", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/", h.handleCreateProduct)
		r.Get("/", h.handleGetProducts)
		r.Get("/{productId}", h.handleGetProduct)
		r.Delete("/{productId}", h.handleDeleteProduct)
		r.Patch("/{productId}", h.handleUpdateProduct)
		r.Patch("/{productId}", h.handleUpdateProduct)
	})
}

func (h handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := chi.URLParam(r, "productId")

	var body dto.UpdateProduct
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = h.productService.UpdateProduct(ctx, productID, body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusOK)
}

func (h handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := chi.URLParam(r, "productId")

	err := h.productService.DeleteProduct(ctx, productID)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusOK)
}

func (h handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := chi.URLParam(r, "productId")

	product, err := h.productService.GetProductByID(ctx, productID)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, product)
}

func (h handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	var q = r.URL.Query()
	var search dto.SearchParams

	search.Term = httputil.ReadQueryString(q, "term", "")
	search.Limit = httputil.ReadQueryInt(q, "limit", 10)
	search.Page = httputil.ReadQueryInt(q, "page", 1)

	ctx := r.Context()
	products, err := h.productService.GetProducts(ctx, search)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, products)
}

func (h handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body dto.CreateProduct
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	err = h.productService.CreateProduct(ctx, body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusOK)
}
