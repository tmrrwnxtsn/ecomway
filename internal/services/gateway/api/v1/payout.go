package v1

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type payoutMethodsRequest struct {
	// Идентификатор клиента
	UserID int64 `query:"user_id" example:"11431" validate:"required"`
	// Валюта платежа в соответствии со стандартом ISO 4217
	Currency string `query:"currency" example:"RUB" validate:"required,iso4217"`
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
}

type payoutMethodsResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Массив платежных методов, доступных для вывода средств
	Methods []method `json:"payout_methods" validate:"required"`
}

func (h *Handler) payoutMethods(c *fiber.Ctx) error {
	ctx := context.Background()

	var req payoutMethodsRequest
	if err := c.QueryParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, err)
	}

	methods, err := h.methodService.AvailableMethods(ctx, model.TransactionTypePayout, req.UserID, req.Currency)
	if err != nil {
		return h.internalErrorResponse(c, err)
	}

	resp := &payoutMethodsResponse{
		Success: true,
		Methods: h.methods(methods, req.LangCode),
	}

	return c.JSON(resp)
}

func (h *Handler) payoutCreate(c *fiber.Ctx) error {
	return c.SendString("Payout created")
}

func (h *Handler) payoutConfirm(c *fiber.Ctx) error {
	return c.SendString("Payout confirmed")
}

func (h *Handler) payoutResendCode(c *fiber.Ctx) error {
	return c.SendString("Payout code resend")
}
