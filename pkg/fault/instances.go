package fault

import (
	"encoding/json"
	"net/http"
)

func NewHTTPResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err, ok := err.(Fault); ok {
		w.WriteHeader(err.StatusCode())
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(New(http.StatusInternalServerError,
		INTERNAL_SERVER_ERROR,
		err.Error(),
		err,
	))
}

func NewConflict(msg string, err error) Fault {
	return New(http.StatusConflict, CONFLICT, msg, err)
}
func NewUnprocessableEntity(msg string, err error) Fault {
	return New(http.StatusUnprocessableEntity, UNPROCESSABLE_ENTITY, msg, err)
}

func NewBadRequest(msg string, err error) Fault {
	return New(http.StatusBadRequest, BAD_REQUEST, msg, err)
}

func NewNotFound(msg string, err error) Fault {
	return New(http.StatusNotFound, NOT_FOUND, msg, err)
}

func NewUnauthorized(msg string, err error) Fault {
	return New(http.StatusUnauthorized, UNAUTHORIZED, msg, err)
}
