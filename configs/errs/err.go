package errs

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string   `json:"message"`
	ErrType string   `json:"error"`
	Code    int      `json:"code"`
	Causes  []*Cause `json:"causes,omitempty"`
}

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

func NewValidationError(message string, causes []*Cause) *Error {
	return NewError(message, "validation_error", http.StatusBadRequest, causes...)
}

func NewBadRequest(message string) *Error {
	return NewError(message, "bad_request", http.StatusBadRequest)
}

func NewNotFound(message string) *Error {
	return NewError(message, "not_found", http.StatusNotFound)
}

func NewInternalServerError(message string) *Error {
	return NewError(message, "internal_server_error", http.StatusInternalServerError)
}

func NewUnauthorized(message string) *Error {
	return NewError(message, "unauthorized", http.StatusUnauthorized)
}

func NewConflict(message string) *Error {
	return NewError(message, "conflict", http.StatusConflict)
}

func (r *Error) ToJSON() string {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func (r *Error) AddCause(field, message string) {
	r.Causes = append(r.Causes, &Cause{Field: field, Message: message})
}

func (r *Error) IsType(errType string) bool {
	return r.ErrType == errType
}
