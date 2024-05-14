package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/translate"
)

type paymentMethodsRequest struct {
	// Идентификатор клиента
	UserID int64 `query:"user_id" example:"1" validate:"required"`
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

// paymentMethods godoc
//
//	@Summary	Получить список способов для пополнения баланса
//	@Tags		Платежи
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		user_id		query		int						true	"Идентификатор клиента"
//	@Param		currency	query		string					true	"Валюта платежа в соответствии со стандартом ISO 4217"
//	@Param		lang_code	query		string					true	"Код языка, обозначение по RFC 5646"
//	@Success	200			{object}	paymentMethodsResponse	"Успешный ответ"
//	@Failure	default		{object}	errorResponse			"Ответ с ошибкой"
//	@Router		/payment/methods [get]
func (h *Handler) paymentMethods(c *fiber.Ctx) error {
	ctx := c.Context()

	var req paymentMethodsRequest
	if err := c.QueryParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	methods, err := h.methodService.AvailableMethods(ctx, model.OperationTypePayment, req.UserID, req.Currency)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	toolsGrouped, err := h.toolService.AvailableToolsGroupedByMethod(ctx, req.UserID)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &paymentMethodsResponse{
		Success: true,
		Methods: h.methods(methods, toolsGrouped, req.LangCode),
	}

	return c.JSON(resp)
}

type paymentReturnURLs struct {
	// URL для возврата клиента, используемый когда результат платежа неизвестен или по умолчанию
	Common string `json:"common" example:"https://example.com" validate:"required"`
	// URL для возврата клиента, используемый при успешном осуществлении платежа
	Success *string `json:"success" example:"https://example.com/success"`
	// URL для возврата клиента, используемый при неуспешном осуществлении платежа
	Fail *string `json:"fail" example:"https://example.com/failed"`
}

type paymentCreateRequest struct {
	// Идентификатор клиента
	UserID int64 `json:"user_id" example:"1" validate:"required"`
	// Идентификатор сохраненного платежного средства
	ToolID string `json:"tool_id" example:"2dc32aa0-000f-5000-8000-16d7bc6cd09f"`
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
	// Объект, содержащий ссылки для возврата клиента для каждого из возможных результатов проведения платежа
	ReturnURLs paymentReturnURLs `json:"return_urls" validate:"required"`
	// Дополнительная информация, специфичная для платежной системы, к которой направляется целевой запрос
	AdditionalData map[string]any `json:"additional_data" swaggertype:"object,string" example:"ip:127.0.0.1,phone_number:+71234567890"`
}

const (
	paymentCreateResponseTypeRedirect = "redirect"
	paymentCreateResponseTypeMessage  = "message"
)

type paymentCreateResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Идентификатор созданного платежа
	OperationID int64 `json:"operation_id" example:"1" validate:"required"`
	// Тип ответа:
	// * Перенаправление клиента на платежную страницу - "redirect"
	// * Текстовое сообщение - "message"
	Type string `json:"type" example:"redirect" validate:"required"`
	// URL платежной страницы, на которую необходимо перенаправить клиента
	RedirectURL string `json:"redirect_url,omitempty" example:"https://securepayments.example.com"`
	// Сообщение, которое необходимо показать клиенту
	Message string `json:"message,omitempty" example:"Баланс пополнен!"`
}

// paymentCreate godoc
//
//	@Summary	Создать запрос на пополнение баланса
//	@Tags		Платежи
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		input	body		paymentCreateRequest	true	"Тело запроса"
//	@Success	200		{object}	paymentCreateResponse	"Успешный ответ"
//	@Failure	default	{object}	errorResponse			"Ответ с ошибкой"
//	@Router		/payment/create [post]
func (h *Handler) paymentCreate(c *fiber.Ctx) error {
	ctx := c.Context()

	var req paymentCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	data := model.CreatePaymentData{
		ReturnURLs:     returnURLsModelFromRequest(req.ReturnURLs),
		AdditionalData: req.AdditionalData,
		ExternalSystem: req.ExternalSystem,
		ExternalMethod: req.ExternalMethod,
		Currency:       req.Currency,
		LangCode:       req.LangCode,
		UserID:         req.UserID,
		ToolID:         req.ToolID,
		Amount:         req.Amount,
	}

	result, err := h.paymentService.Create(ctx, data)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &paymentCreateResponse{
		Success:     true,
		OperationID: result.OperationID,
	}

	switch result.Status {
	case model.OperationStatusSuccess, model.OperationStatusFailed:
		resp.Type = paymentCreateResponseTypeMessage
		resp.Message = h.paymentMessageFromResult(req.LangCode, result.Status)
	default:
		resp.Type = paymentCreateResponseTypeRedirect
		resp.RedirectURL = result.RedirectURL
	}

	return c.JSON(resp)
}

func returnURLsModelFromRequest(returnURLs paymentReturnURLs) model.ReturnURLs {
	result := model.ReturnURLs{
		Common: returnURLs.Common,
	}

	if returnURLs.Success != nil {
		result.Success = *returnURLs.Success
	}

	if returnURLs.Fail != nil {
		result.Fail = *returnURLs.Fail
	}

	return result
}

func (h *Handler) paymentMessageFromResult(langCode string, resultStatus model.OperationStatus) string {
	var translateKey string

	switch resultStatus {
	case model.OperationStatusSuccess:
		translateKey = translate.KeyPaymentSuccessful
	case model.OperationStatusFailed:
		translateKey = translate.KeyPaymentRejected
	}

	return h.translator.Translate(langCode, translateKey)
}
