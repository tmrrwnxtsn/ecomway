package v1

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/translate"
)

type payoutMethodsRequest struct {
	// Идентификатор клиента
	UserID string `query:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `query:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
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

// payoutMethods godoc
//
//	@Summary	Получить список способов для вывода средств
//	@Tags		Выплаты
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		user_id		query		string					true	"Идентификатор клиента"
//	@Param		session_id	query		string					true	"Идентификатор сессии клиента"
//	@Param		currency	query		string					true	"Валюта выплаты в соответствии со стандартом ISO 4217"
//	@Param		lang_code	query		string					true	"Код языка, обозначение по RFC 5646"
//	@Success	200			{object}	payoutMethodsResponse	"Успешный ответ"
//	@Failure	default		{object}	errorResponse			"Ответ с ошибкой"
//	@Router		/payout/methods [get]
func (h *Handler) payoutMethods(c *fiber.Ctx) error {
	ctx := c.Context()

	var req payoutMethodsRequest
	if err := c.QueryParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	methods, err := h.methodService.AvailableMethods(ctx, model.OperationTypePayout, req.UserID, req.Currency)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	toolsGrouped, err := h.toolService.AvailableToolsGroupedByMethod(ctx, req.UserID)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &payoutMethodsResponse{
		Success: true,
		Methods: h.methods(methods, toolsGrouped, req.LangCode),
	}

	return c.JSON(resp)
}

type payoutCreateRequest struct {
	// Идентификатор клиента
	UserID string `json:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `json:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
	// Идентификатор сохраненного платежного средства
	ToolID string `json:"tool_id" example:"2dc32aa0-000f-5000-8000-16d7bc6cd09f" validate:"required"`
	// Сумма выплаты в минорных единицах валюты (копейки, центы и т.п.)
	Amount int64 `json:"amount" example:"10000" validate:"required,gte=100"`
	// Валюта выплаты в соответствии со стандартом ISO 4217
	Currency string `json:"currency" example:"RUB" validate:"required,iso4217"`
	// Внутренний код платежной системы, к которой направляется целевой запрос
	ExternalSystem string `json:"external_system" example:"yookassa" validate:"required"`
	// Внутренний код платежного метода платежной системы, к которой направляется целевой запрос
	ExternalMethod string `json:"external_method" example:"yookassa_bank_card" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `json:"lang_code" example:"en" validate:"required"`
	// Дополнительная информация, специфичная для платежной системы, к которой направляется целевой запрос
	AdditionalData map[string]any `json:"additional_data" swaggertype:"object,string" example:"ip:127.0.0.1,phone_number:+71234567890"`
}

type payoutCreateResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Идентификатор созданной выплаты
	OperationID int64 `json:"operation_id" example:"1" validate:"required"`
}

// payoutCreate godoc
//
//	@Summary	Создать запрос на вывод средств
//	@Tags		Выплаты
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		input	body		payoutCreateRequest		true	"Тело запроса"
//	@Success	200		{object}	payoutCreateResponse	"Успешный ответ"
//	@Failure	default	{object}	errorResponse			"Ответ с ошибкой"
//	@Router		/payout/create [post]
func (h *Handler) payoutCreate(c *fiber.Ctx) error {
	ctx := c.Context()

	var req payoutCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	data := model.CreatePayoutData{
		AdditionalData: req.AdditionalData,
		ExternalSystem: req.ExternalSystem,
		ExternalMethod: req.ExternalMethod,
		Currency:       req.Currency,
		LangCode:       req.LangCode,
		UserID:         req.UserID,
		ToolID:         req.ToolID,
		Amount:         req.Amount,
	}

	result, err := h.payoutService.Create(ctx, data)
	if err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal {
				switch perr.Code {
				case perror.CodeObjectNotFound:
					return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
				case perror.CodeToolHasBeenRemoved:
					return h.forbiddenOnRemovedToolErrorResponse(c, req.LangCode, perr)
				}
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &payoutCreateResponse{
		Success:     true,
		OperationID: result.OperationID,
	}

	return c.JSON(resp)
}

type payoutConfirmRequest struct {
	// Идентификатор клиента
	UserID string `json:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `json:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `json:"lang_code" example:"en" validate:"required"`
	// Код подтверждения вывода средств
	ConfirmationCode string `json:"confirmation_code" example:"123456" validate:"required"`
}

type payoutConfirmResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Сообщение, которое необходимо показать клиенту
	Message string `json:"message" example:"Вывод средств подтвержден." validate:"required"`
}

// payoutConfirm godoc
//
//	@Summary	Подтвердить запрос на вывод средств
//	@Tags		Выплаты
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id		path		int						true	"Идентификатор операции"
//	@Param		input	body		payoutConfirmRequest	true	"Тело запроса"
//	@Success	200		{object}	payoutConfirmResponse	"Успешный ответ"
//	@Failure	default	{object}	errorResponse			"Ответ с ошибкой"
//	@Router		/payout/{id}/confirm [put]
func (h *Handler) payoutConfirm(c *fiber.Ctx) error {
	ctx := c.Context()

	var req payoutConfirmRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	opID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("failed to parse id as int: %w", err)
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	confirmData := model.ConfirmPayoutData{
		ConfirmationCode: req.ConfirmationCode,
		LangCode:         req.LangCode,
		UserID:           req.UserID,
		OperationID:      opID,
	}

	if err = h.payoutService.Confirm(ctx, confirmData); err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal {
				switch perr.Code {
				case perror.CodeObjectNotFound:
					return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
				case perror.CodeToolHasBeenRemoved:
					return h.forbiddenOnRemovedToolErrorResponse(c, req.LangCode, perr)
				case perror.CodeWrongConfirmationCode:
					return h.wrongConfirmationCodeErrorResponse(c, req.LangCode, perr)
				case perror.CodeConfirmationAttemptsExceeded:
					return h.wrongCodeLimitExceededErrorResponse(c, req.LangCode, perr)
				}
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &payoutConfirmResponse{
		Success: true,
		Message: h.translator.Translate(req.LangCode, translate.KeyPayoutConfirmed),
	}

	return c.JSON(resp)
}

type payoutResendCodeRequest struct {
	// Идентификатор клиента
	UserID string `json:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `json:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `json:"lang_code" example:"en" validate:"required"`
}

type payoutResendCodeResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
}

// payoutResendCode godoc
//
//	@Summary	Отправить код подтверждения выплаты повторно
//	@Tags		Выплаты
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id		path		int							true	"Идентификатор операции"
//	@Param		input	body		payoutResendCodeRequest		true	"Тело запроса"
//	@Success	200		{object}	payoutResendCodeResponse	"Успешный ответ"
//	@Failure	default	{object}	errorResponse				"Ответ с ошибкой"
//	@Router		/payout/{id}/resend-code [put]
func (h *Handler) payoutResendCode(c *fiber.Ctx) error {
	ctx := c.Context()

	var req payoutResendCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	opID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("failed to parse id as int: %w", err)
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err = h.payoutService.ResendCode(ctx, opID, req.UserID, req.LangCode); err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal && perr.Code == perror.CodeObjectNotFound {
				return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &payoutResendCodeResponse{
		Success: true,
	}

	return c.JSON(resp)
}
