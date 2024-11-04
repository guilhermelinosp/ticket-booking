package middlewares

import (
	"strings"
	"ticket-booking/configs/errs"
	"ticket-booking/configs/logs"
	"ticket-booking/services"

	"github.com/gofiber/fiber/v2"
)

func Auth(tokenization services.Tokenization) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			logs.Error("Middleware.Auth: Missing Authorization header", nil)
			return errs.NewUnauthorized(ctx, "Missing Authorization header")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logs.Error("Middleware.Auth: Invalid Authorization header format", nil)
			return errs.NewUnauthorized(ctx, "Invalid Authorization header format")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		valid, err := tokenization.ValidateToken(token)
		if err != nil || !valid {
			logs.Error("Middleware.Auth: Invalid or expired token", err)
			return errs.NewUnauthorized(ctx, "Invalid or expired token")
		}

		return ctx.Next()
	}
}
