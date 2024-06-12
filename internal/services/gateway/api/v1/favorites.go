package v1

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type favoritesRequest struct {
	// Идентификатор клиента
	UserID int64 `json:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `json:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
	// Внутренний код платежной системы, платежный метод которой необходимо занести в "Избранное"
	ExternalSystem string `json:"external_system" example:"yookassa" validate:"required"`
	// Внутренний код платежного метода платежной системы, которого необходимо занести в "Избранное"
	ExternalMethod string `json:"external_method" example:"yookassa_bank_card" validate:"required"`
	// Валюта платежного метода платежной системы, которого необходимо занести в "Избранное"
	Currency string `json:"currency" example:"RUB" validate:"required,iso4217"`
	// Код языка, обозначение по RFC 5646
	LangCode string `json:"lang_code" example:"en" validate:"required"`
}

type favoritesResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
}

// favoritesAdd godoc
//
//	@Summary	Добавить платежную систему в избранные способы оплаты (выплаты)
//	@Tags		Избранное
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		operation_type	path		string				true	"Тип транзакции"
//	@Param		input			body		favoritesRequest	true	"Тело запроса"
//	@Success	200				{object}	favoritesResponse	"Успешный ответ"
//	@Failure	default			{object}	errorResponse		"Ответ с ошибкой"
//	@Router		/favorites/{operation_type} [post]
func (h *Handler) favoritesAdd(c *fiber.Ctx) error {
	ctx := c.Context()

	var req favoritesRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	data := model.FavoritesData{
		OperationType:  model.OperationType(c.Params("operation_type")),
		Currency:       req.Currency,
		ExternalSystem: req.ExternalSystem,
		ExternalMethod: req.ExternalMethod,
		UserID:         req.UserID,
	}

	if err := h.favoritesService.Add(ctx, data); err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal && perr.Code == perror.CodeObjectNotFound {
				return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	return c.JSON(&favoritesResponse{
		Success: true,
	})
}

// favoritesRemove godoc
//
//	@Summary	Удалить платежную систему из избранных способов оплаты (выплаты)
//	@Tags		Избранное
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		operation_type	path		string				true	"Тип транзакции"
//	@Param		input			body		favoritesRequest	true	"Тело запроса"
//	@Success	200				{object}	favoritesResponse	"Успешный ответ"
//	@Failure	default			{object}	errorResponse		"Ответ с ошибкой"
//	@Router		/favorites/{operation_type} [delete]
func (h *Handler) favoritesRemove(c *fiber.Ctx) error {
	ctx := c.Context()

	var req favoritesRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	data := model.FavoritesData{
		OperationType:  model.OperationType(c.Params("operation_type")),
		Currency:       req.Currency,
		ExternalSystem: req.ExternalSystem,
		ExternalMethod: req.ExternalMethod,
		UserID:         req.UserID,
	}

	if err := h.favoritesService.Remove(ctx, data); err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal && perr.Code == perror.CodeObjectNotFound {
				return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	return c.JSON(&favoritesResponse{
		Success: true,
	})
}
