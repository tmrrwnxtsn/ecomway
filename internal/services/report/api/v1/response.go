package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/translate"
)

const (
	errorCodeInvalidRequest = "InvalidRequest"
	errorCodeInvalidAPIKey  = "InvalidAPIKey"
	errorCodeInternalError  = "InternalError"
)

type errorContent struct {
	// Код ошибки
	Code string `json:"code" example:"InvalidRequest" validate:"required"`
	// Описание ошибки для разработки
	Description string `json:"description" example:"user_id param is required" validate:"required"`
	// Сообщение об ошибке для клиента
	Message string `json:"message,omitempty" example:"Internal server error occurred. Please try again later." validate:"optional"`
}

type errorResponse struct {
	// Результат обработки запроса (всегда false)
	Success bool `json:"success" example:"false" validate:"required"`
	// Развернутая информация об ошибке
	Error errorContent `json:"error" validate:"required"`
}

func (h *Handler) requestValidationErrorResponse(c *fiber.Ctx, langCode string, err error) error {
	return c.Status(http.StatusBadRequest).JSON(&errorResponse{
		Success: false,
		Error: errorContent{
			Code:        errorCodeInvalidRequest,
			Description: err.Error(),
			Message:     h.translator.Translate(langCode, translate.KeyUnexpectedError),
		},
	})
}

func (h *Handler) internalErrorResponse(c *fiber.Ctx, langCode string, err error) error {
	return c.Status(http.StatusInternalServerError).JSON(&errorResponse{
		Success: false,
		Error: errorContent{
			Code:        errorCodeInternalError,
			Description: err.Error(),
			Message:     h.translator.Translate(langCode, translate.KeyUnexpectedError),
		},
	})
}
