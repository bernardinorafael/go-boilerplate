package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const maxRequestBodyBytes = 1_048_576 // 1MB

// WriteSuccessResponse writes a JSON success response with the specified HTTP status code.
// The JSON response will contain a single key "message" with the value "success".
func WriteSuccess(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

// WriteJSONResponse writes a JSON response with the specified HTTP status code and the provided data.
// The data to be encoded as JSON should be passed as the 'dst' parameter.
func WriteJSON(w http.ResponseWriter, code int, dst any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(dst)
}

// ReadQueryInt reads a query string parameter from the URL values and parses it into an integer.
// If the parameter is missing or cannot be parsed, the provided default value 'defval' is returned.
//
// Example:
//
//	page := httputil.ReadQueryInt(r.URL.Query(), "page", 1)
//
// If the query string parameter "page" is not present, the default value 1 is returned.
// If the query string parameter "page" is present but cannot be parsed as an integer,
// the default value 1 is returned.
func ReadQueryInt(qs url.Values, key string) int {
	val := qs.Get(key)
	if val == "" {
		return 0
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}

	return i
}

// ReadQueryBool reads a query string param from the URL and returns it as a boolean
// If the parameter is missing or cannot be parsed, the provided default value 'defval' is returned.
//
// Example:
//
//	includeArchived := httputil.ReadQueryBool(r.URL.Query(), "include_archived", false)
//
// If the query string parameter "include_archived" is not present, the default value false is returned.
// If the query string parameter "include_archived" is present but cannot be parsed as a boolean,
func ReadQueryBool(qs url.Values, key string) bool {
	val := qs.Get(key)
	if val == "" {
		return false
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return b
}

// ReadQueryString reads a query string parameter from the URL values and returns it as a string.
// If the parameter is missing, the provided default value 'defval' is returned.
//
// Example:
//
//	sort := httputil.ReadQueryString(r.URL.Query(), "sort", "asc")
//
// If the query string parameter "sort" is not present, the default value "asc" is returned.
func ReadQueryString(qs url.Values, key string) string {
	val := qs.Get(key)
	if val == "" {
		return ""
	}
	return val
}

// ReadRequestBody reads and parses the JSON body of an HTTP request into the provided destination struct.
// It limits the size of the request body to 1MB and returns detailed error messages for various parsing issues.
//
// Example:
//
//	var body struct {
//		ID string `json:"id"`
//		Name string `json:"name"`
//	}
//
//	err := httputil.ReadRequestBody(w, r, &body)
//	if err != nil {
//		// handle error here
//	}
func ReadRequestBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxRequestBodyBytes))

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err := d.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			// JSON syntax error in the request body
			// Offset is the exact byte where the error occurred
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			// JSON value and struct type do not match
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			// io.EOF (End of File) indicates that there are no more bytes left to read
			return ErrEmptyRequestBody
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return ErrUnknownRequestBodyKey
		case errors.As(err, &invalidUnmarshalError):
			// Received a non-nil pointer into Decode()
			panic(err)
		default:
			return err
		}
	}

	// Calling decode again to check if there's more data after the JSON object
	// This will return an io.EOF error, indicating that the client sent more data
	err = d.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
