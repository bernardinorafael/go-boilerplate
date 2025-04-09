package fault

import (
	"encoding/json"
	"net/http"
)

// NewHTTPError receives an error and writes it to the response writer
// It sets the content type to application/json and writes the error
// If the error is not a Fault, it writes a new InternalServerError
func NewHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err, ok := err.(*Fault); ok {
		w.WriteHeader(err.GetHTTPCode())
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(New(
		"an unexpected error occurred",
		WithHTTPCode(http.StatusInternalServerError),
		WithTag(INTERNAL_SERVER_ERROR),
		WithError(err),
	))
}

func NewBadRequest(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusBadRequest),
		WithTag(BAD_REQUEST),
	)
}

func NewNotFound(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusNotFound),
		WithTag(NOT_FOUND),
	)
}

func NewInternalServerError(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusInternalServerError),
		WithTag(INTERNAL_SERVER_ERROR),
	)
}

func NewUnauthorized(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusUnauthorized),
		WithTag(UNAUTHORIZED),
	)
}

func NewForbidden(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusForbidden),
		WithTag(FORBIDDEN),
	)
}

func NewConflict(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusConflict),
		WithTag(CONFLICT),
	)
}

func NewTooManyRequests(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusTooManyRequests),
		WithTag(TOO_MANY_REQUESTS),
	)
}

func NewUnprocessableEntity(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusUnprocessableEntity),
		WithTag(UNPROCESSABLE_ENTITY),
	)
}
