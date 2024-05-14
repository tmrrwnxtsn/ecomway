package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewAccessLog() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		stop := time.Now()

		slog.Info("request has been processed",
			"request_method", c.Method(),
			"request_route", c.Path(),
			"request_body", string(c.Request().Body()),
			"response_status", c.Response().StatusCode(),
			"response_latency", stop.Sub(start).Milliseconds(),
			"response_body", string(c.Response().Body()),
		)

		return err
	}
}
