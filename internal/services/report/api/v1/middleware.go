package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/translate"
)

func (h *Handler) authorizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get(fiber.HeaderAuthorization) != h.apiKey {
			return h.authorizationFailedErrorResponse(c)
		}
		return c.Next()
	}
}

func (h *Handler) authorizationFailedErrorResponse(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(&errorResponse{
		Success: false,
		Error: errorContent{
			Code:        errorCodeInvalidAPIKey,
			Description: "Authorization header does not contain valid API Key",
			Message:     h.translator.Translate("", translate.KeyInvalidAPIKey),
		},
	})
}
