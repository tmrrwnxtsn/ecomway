package v1

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type paymentMethodsResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Массив платежных методов, доступных для пополнения баланса
	Methods []method `json:"payment_methods" validate:"required"`
}

func (h *Handler) paymentMethods(c *fiber.Ctx) error {
	ctx := context.Background()

	var req request
	if err := c.QueryParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, err)
	}

	methods, err := h.methodService.AvailableMethods(ctx, req.UserID, model.TransactionTypePayment)
	if err != nil {
		return h.internalErrorResponse(c, err)
	}

	resp := &paymentMethodsResponse{
		Success: true,
		Methods: h.methods(methods, req.LangCode),
	}

	return c.JSON(resp)
}

func (h *Handler) paymentCreate(c *fiber.Ctx) error {
	return c.SendString("Payment created")
}
