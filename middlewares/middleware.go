package middlewares

import (
	"time"

	"ticket-booking/configs/logs"

	"github.com/gofiber/fiber/v2"
)

// Logger logs the request information.
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		logs.Request(c.Method(), c.Path(), c.Response().StatusCode(), duration.String())

		return err
	}
}
