package middleware

import (
	"net/http"
	"strconv"

	"github.com/bernardinorafael/go-boilerplate/pkg/metric"
)

// WithMetrics is a middleware that records the metrics of the HTTP requests
func withMetrics(metrics *metric.Metric) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			// After the request is processed, record the metrics
			status := strconv.Itoa(rw.statusCode)
			metrics.RecordHTTPRequest(r.Method, r.URL.Path, status)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
