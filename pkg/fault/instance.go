package fault

import (
	"encoding/json"
	"errors"
	"net/http"
)

// NewHTTPError receives an error and writes it to the response writer
// It sets the content type to application/json and writes the error
// If the error is not a Fault, it writes a new InternalServerError
func NewHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var fault *Fault
	if errors.As(err, &fault) {
		w.WriteHeader(fault.GetHTTPCode())
		_ = json.NewEncoder(w).Encode(fault)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(
		New(
			"an unexpected error occurred",
			WithHTTPCode(http.StatusInternalServerError),
			WithTag(InternalServerError),
			WithError(err),
		),
	)
}

func NewBadRequest(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusBadRequest),
		WithTag(BadRequest),
	)
}

func NewNotFound(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusNotFound),
		WithTag(NotFound),
	)
}

func NewInternalServerError(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusInternalServerError),
		WithTag(InternalServerError),
	)
}

func NewUnauthorized(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusUnauthorized),
		WithTag(Unauthorized),
	)
}

func NewForbidden(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusForbidden),
		WithTag(Forbidden),
	)
}

func NewConflict(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusConflict),
		WithTag(Conflict),
	)
}

func NewTooManyRequests(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusTooManyRequests),
		WithTag(TooManyRequests),
	)
}

func NewUnprocessableEntity(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusUnprocessableEntity),
		WithTag(UnprocessableEntity),
	)
}
