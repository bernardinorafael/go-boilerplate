package fault

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		options  []func(*Fault)
		expected *Fault
	}{
		{
			name:    "basic fault",
			message: "test error",
			expected: &Fault{
				Tag:        Untagged,
				Message:    "test error",
				HTTPCode:   http.StatusBadRequest,
				FieldError: []FieldError{},
			},
		},
		{
			name:    "fault with tag",
			message: "not found",
			options: []func(*Fault){WithTag(NotFound)},
			expected: &Fault{
				Tag:        NotFound,
				Message:    "not found",
				HTTPCode:   http.StatusBadRequest,
				FieldError: []FieldError{},
			},
		},
		{
			name:    "fault with HTTP code",
			message: "unauthorized",
			options: []func(*Fault){WithHTTPCode(http.StatusUnauthorized)},
			expected: &Fault{
				Tag:        Untagged,
				Message:    "unauthorized",
				HTTPCode:   http.StatusUnauthorized,
				FieldError: []FieldError{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.message, tt.options...)
			if result.Tag != tt.expected.Tag {
				t.Errorf("expected tag %v, got %v", tt.expected.Tag, result.Tag)
			}
			if result.Message != tt.expected.Message {
				t.Errorf("expected message %v, got %v", tt.expected.Message, result.Message)
			}
			if result.HTTPCode != tt.expected.HTTPCode {
				t.Errorf("expected HTTP code %v, got %v", tt.expected.HTTPCode, result.HTTPCode)
			}
		})
	}
}

func TestWithValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected []FieldError
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: []FieldError{},
		},
		{
			name: "valid validation error",
			err:  errors.New("email: invalid email format, age: must be positive"),
			expected: []FieldError{
				{Field: "email", Message: "invalid email format"},
				{Field: "age", Message: "must be positive"},
			},
		},
		{
			name: "malformed validation error",
			err:  errors.New("invalid format without colon"),
			expected: []FieldError{
				{Field: "general", Message: "invalid format without colon"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fault := New("test")
			WithValidationError(tt.err)(fault)

			if len(fault.FieldError) != len(tt.expected) {
				t.Errorf("expected %d field errors, got %d", len(tt.expected), len(fault.FieldError))
			}

			for i, expected := range tt.expected {
				if i >= len(fault.FieldError) {
					break
				}
				if fault.FieldError[i].Field != expected.Field {
					t.Errorf("expected field %v, got %v", expected.Field, fault.FieldError[i].Field)
				}
				if fault.FieldError[i].Message != expected.Message {
					t.Errorf("expected message %v, got %v", expected.Message, fault.FieldError[i].Message)
				}
			}
		})
	}
}

func TestNewHTTPError(t *testing.T) {
	t.Run("fault error", func(t *testing.T) {
		w := httptest.NewRecorder()
		fault := New("test error", WithHTTPCode(http.StatusNotFound))
		NewHTTPError(w, fault)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected content type %v, got %v", "application/json", contentType)
		}
	})

	t.Run("non-fault error", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := errors.New("some error")
		NewHTTPError(w, err)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestGetTag(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected Tag
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: Untagged,
		},
		{
			name:     "fault error",
			err:      New("test", WithTag(NotFound)),
			expected: NotFound,
		},
		{
			name:     "wrapped fault error",
			err:      errors.New("wrapper: " + New("test", WithTag(BadRequest)).Error()),
			expected: Untagged,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTag(tt.err)
			if result != tt.expected {
				t.Errorf("expected tag %v, got %v", tt.expected, result)
			}
		})
	}
}
