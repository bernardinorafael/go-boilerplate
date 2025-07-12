package fault

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Fault struct {
	Tag        Tag          `json:"tag"`
	Message    string       `json:"message"`
	FieldError []FieldError `json:"fields"`

	HTTPCode int   `json:"-"`
	Err      error `json:"-"`
}

// New instantiates a new Fault with the given message
// The message is used to describe the error in detail
//
// The default HTTP code is 400.
func New(msg string, options ...func(*Fault)) *Fault {
	var validations = make([]FieldError, 0)

	fault := Fault{
		Err:        nil,
		Tag:        Untagged,
		HTTPCode:   http.StatusBadRequest,
		Message:    msg,
		FieldError: validations,
	}

	for _, fn := range options {
		fn(&fault)
	}

	return &fault
}

func WithValidationError(err error) func(*Fault) {
	if err == nil {
		return func(f *Fault) {}
	}

	var validations []FieldError
	splittedError := strings.SplitSeq(err.Error(), ";")

	for validation := range splittedError {
		validation = strings.TrimSpace(validation)
		if validation == "" {
			continue
		}

		parts := strings.SplitN(validation, ":", 2)
		if len(parts) != 2 {
			validations = append(validations, FieldError{
				Field:   "general",
				Message: validation,
			})
			continue
		}

		field := strings.TrimSpace(parts[0])
		msg := strings.TrimSpace(parts[1])

		// If the field or message is empty, skip it
		if field == "" || msg == "" {
			continue
		}

		validations = append(validations, FieldError{
			Field:   field,
			Message: msg,
		})
	}

	return func(f *Fault) {
		f.FieldError = validations
	}
}

// WithHTTPCode sets the HTTP code for the fault
func WithHTTPCode(code int) func(*Fault) {
	return func(f *Fault) {
		f.HTTPCode = code
	}
}

// WithError sets the error for the fault
func WithError(err error) func(*Fault) {
	return func(f *Fault) {
		if err == nil {
			return
		}
		f.Err = err
	}
}

// WithTag sets the tag for the fault
func WithTag(tag Tag) func(*Fault) {
	return func(f *Fault) {
		f.Tag = tag
	}
}

// WithFieldError sets the field errors for the fault
func WithFieldError(args ...FieldError) func(*Fault) {
	return func(f *Fault) {
		f.FieldError = args
	}
}

// GetHTTPCode returns the HTTP code for the fault
func (f *Fault) GetHTTPCode() int {
	return f.HTTPCode
}

func (f *Fault) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", f.Tag, f.Message, f.Err)
	}
	return fmt.Sprintf("%s: %s", f.Tag, f.Message)
}

func (f *Fault) Is(target error) bool {
	var fault *Fault
	return errors.As(target, &fault)
}

func (f *Fault) Unwrap() error {
	return f.Err
}
