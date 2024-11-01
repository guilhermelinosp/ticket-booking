package errs

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Error represents an error response.
type Error struct {
	Status  int    `json:"status"`
	ErrType string `json:"error"`
	Message string `json:"message"`
}

// NewError creates a new Error instance with optional causes.
func NewError(message, errType string, status int) *Error {
	return &Error{
		Status:  status,
		ErrType: errType,
		Message: message,
	}
}

// NewBadRequest creates a bad request error response.
func NewBadRequest(ctx *fiber.Ctx, message string) error {
	err := NewError(message, "bad_request_error", http.StatusBadRequest)
	return ctx.Status(http.StatusBadRequest).JSON(err)
}

// NewNotFound creates a not found error response.
func NewNotFound(ctx *fiber.Ctx, message string) error {
	err := NewError(message, "not_found_error", http.StatusNotFound)
	return ctx.Status(http.StatusNotFound).JSON(err)
}

// NewInternalServerError creates an internal server error response.
func NewInternalServerError(ctx *fiber.Ctx, message string) error {
	err := NewError(message, "internal_server_error", http.StatusInternalServerError)
	return ctx.Status(http.StatusInternalServerError).JSON(err)
}

// NewUnauthorized creates an unauthorized error response.
func NewUnauthorized(ctx *fiber.Ctx, message string) error {
	err := NewError(message, "unauthorized_error", http.StatusUnauthorized)
	return ctx.Status(http.StatusUnauthorized).JSON(err)
}
