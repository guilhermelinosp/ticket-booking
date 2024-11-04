package errs

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Status  int    `json:"status"`
	ErrType string `json:"error"`
	Message string `json:"message"`
}

func NewError(message, errType string, status int) *Error {
	return &Error{
		Status:  status,
		ErrType: errType,
		Message: message,
	}
}

func NewBadRequest(ctx *fiber.Ctx, message string) error {
	err := NewError(message, "bad_request_error", http.StatusBadRequest)
	return ctx.Status(http.StatusBadRequest).JSON(err)
}

func NewNotFound(ctx *fiber.Ctx, message string) error {
	err := NewError(message, "not_found_error", http.StatusNotFound)
	return ctx.Status(http.StatusNotFound).JSON(err)
}

func NewInternalServerError(ctx *fiber.Ctx, message string) error {
	err := NewError(message, "internal_server_error", http.StatusInternalServerError)
	return ctx.Status(http.StatusInternalServerError).JSON(err)
}

func NewUnauthorized(ctx *fiber.Ctx, message string) error {
	err := NewError(message, "unauthorized_error", http.StatusUnauthorized)
	return ctx.Status(http.StatusUnauthorized).JSON(err)
}
