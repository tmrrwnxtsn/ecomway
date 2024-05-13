package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
)

const (
	errorCodeInvalidRequest = "InvalidRequest"
	errorCodeInvalidAPIKey  = "InvalidAPIKey"
	errorCodeInternalError  = "InternalError"
	errorCodeObjectNotFound = "ObjectNotFound"
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

func (h *Handler) requestValidationErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusBadRequest).JSON(&errorResponse{
		Success: false,
		Error: errorContent{
			Code:        errorCodeInvalidRequest,
			Description: err.Error(),
		},
	})
}

func (h *Handler) internalErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusInternalServerError).JSON(&errorResponse{
		Success: false,
		Error: errorContent{
			Code:        errorCodeInternalError,
			Description: err.Error(),
		},
	})
}

func (h *Handler) objectNotFoundErrorResponse(c *fiber.Ctx, perr *perror.Error) error {
	return c.Status(http.StatusNotFound).JSON(&errorResponse{
		Success: false,
		Error: errorContent{
			Code:        errorCodeObjectNotFound,
			Description: perr.Description,
		},
	})
}
