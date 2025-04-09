package fault

import (
	"fmt"
	"net/http"
)

type Fault struct {
	HTTPCode int    `json:"-"`
	Err      error  `json:"-"`
	Tag      Tag    `json:"tag"`
	Message  string `json:"message"`
}

// New instantiates a new Fault with the given message
// The message is used to describe the error in detail
//
// The default HTTP code is 400.
func New(msg string, options ...func(*Fault)) *Fault {
	fault := Fault{
		HTTPCode: http.StatusBadRequest,
		Err:      nil,
		Tag:      UNTAGGED,
		Message:  msg,
	}

	for _, fn := range options {
		fn(&fault)
	}

	return &fault
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

// GetHTTPCode returns the HTTP code for the fault
func (f *Fault) GetHTTPCode() int {
	return f.HTTPCode
}

func (f *Fault) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s:%s (caused by: %v)", f.Tag, f.Message, f.Err)
	}
	return fmt.Sprintf("%s:%s", f.Tag, f.Message)
}

func (f *Fault) Is(target error) bool {
	_, ok := target.(*Fault)
	return ok
}

func (f *Fault) Unwrap() error {
	return f.Err
}
