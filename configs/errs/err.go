package errs

import (
	"encoding/json"
	"net/http"
)

// Error represents an error response.
type Error struct {
	Message string   `json:"message"`
	ErrType string   `json:"error"`
	Code    int      `json:"code"`
	Causes  []*Cause `json:"causes,omitempty"`
}

// Cause represents a cause of an error.
type Cause struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewError creates a new Error instance with optional causes as pointers.
func NewError(message, error string, code int, causes ...*Cause) *Error {
	return &Error{
		Message: message,
		ErrType: error,
		Code:    code,
		Causes:  append([]*Cause{}, causes...), // Initialize as empty slice if no causes are provided
	}
}

// NewValidationError creates a new validation error instance.
func NewValidationError(message string, causes []*Cause) *Error {
	return NewError(message, "validation_error", http.StatusBadRequest, causes...)
}

// NewBadRequest creates a new bad request error instance.
func NewBadRequest(message string) *Error {
	return NewError(message, "bad_request", http.StatusBadRequest)
}

// NewNotFound creates a new not found error instance.
func NewNotFound(message string) *Error {
	return NewError(message, "not_found", http.StatusNotFound)
}

// NewInternalServerError creates a new internal server error instance.
func NewInternalServerError(message string) *Error {
	return NewError(message, "internal_server_error", http.StatusInternalServerError)
}

// NewUnauthorized creates a new unauthorized error instance.
func NewUnauthorized(message string) *Error {
	return NewError(message, "unauthorized", http.StatusUnauthorized)
}

// NewConflict creates a new conflict error instance.
func NewConflict(message string) *Error {
	return NewError(message, "conflict", http.StatusConflict)
}

// ToJSON converts the error instance to a JSON string.
func (r *Error) ToJSON() string {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// AddCause appends a new cause to the error instance.
func (r *Error) AddCause(field, message string) {
	r.Causes = append(r.Causes, &Cause{Field: field, Message: message})
}

// IsType checks if the error type matches the provided type.
func (r *Error) IsType(errType string) bool {
	return r.ErrType == errType
}
