package v1

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type paymentMethodsRequest struct {
	// Идентификатор клиента
	UserID int64 `query:"user_id" example:"11431" validate:"required"`
	// Валюта платежа в соответствии со стандартом ISO 4217
	Currency string `query:"currency" example:"RUB" validate:"required,iso4217"`
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
}

type paymentMethodsResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Массив платежных методов, доступных для пополнения баланса
	Methods []method `json:"payment_methods" validate:"required"`
}

func (h *Handler) paymentMethods(c *fiber.Ctx) error {
	ctx := context.Background()

	var req paymentMethodsRequest
	if err := c.QueryParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, err)
	}

	methods, err := h.methodService.AvailableMethods(ctx, model.OperationTypePayment, req.UserID, req.Currency)
	if err != nil {
		return h.internalErrorResponse(c, err)
	}

	resp := &paymentMethodsResponse{
		Success: true,
		Methods: h.methods(methods, req.LangCode),
	}

	return c.JSON(resp)
}

type paymentCreateRequest struct {
	// Идентификатор клиента
	UserID int64 `json:"user_id" example:"11431" validate:"required"`
	// Сумма платежа в минорных единицах валюты (копейки, центы и т.п.)
	Amount int64 `json:"amount" example:"10000" validate:"required,gte=100"`
	// Валюта платежа в соответствии со стандартом ISO 4217
	Currency string `json:"currency" example:"RUB" validate:"required,iso4217"`
	// Внутренний код платежной системы, к которой направляется целевой запрос
	ExternalSystem string `json:"external_system" example:"yookassa" validate:"required"`
	// Внутренний код платежного метода платежной системы, к которой направляется целевой запрос
	ExternalMethod string `json:"external_method" example:"yookassa_bank_card" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `json:"lang_code" example:"en" validate:"required"`
	// Дополнительная информация, специфичная для платежной системы, к которой направляется целевой запрос
	AdditionalData map[string]any `json:"additional_data" example:"ip:127.0.0.1,phone_number:+71234567890"`
}

type paymentCreateResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// URL платежной страницы, на которую необходимо перенаправить клиента
	RedirectURL string `json:"redirect_url" example:"example.com" validate:"required"`
}

func (h *Handler) paymentCreate(c *fiber.Ctx) error {
	ctx := context.Background()

	var req paymentCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, err)
	}

	data := model.CreatePaymentData{
		AdditionalData: req.AdditionalData,
		ExternalSystem: req.ExternalSystem,
		ExternalMethod: req.ExternalMethod,
		Currency:       req.Currency,
		LangCode:       req.LangCode,
		UserID:         req.UserID,
		Amount:         req.Amount,
	}

	result, err := h.paymentService.Create(ctx, data)
	if err != nil {
		return h.internalErrorResponse(c, err)
	}

	resp := &paymentCreateResponse{
		Success:     true,
		RedirectURL: result.RedirectURL,
	}

	return c.JSON(resp)
}
