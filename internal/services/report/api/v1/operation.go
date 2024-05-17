package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type operation struct {
	// Идентификатор операции
	ID int64 `json:"id" example:"1" validate:"required"`
	// Идентификатор клиента
	UserID int64 `json:"client_id" example:"1" validate:"required"`
	// Тип операции
	Type string `json:"type" example:"payment" validate:"required"`
	// Валюта операции
	Currency string `json:"currency" example:"RUB" validate:"required"`
	// Сумма операции
	Amount float64 `json:"amount" example:"121.01" validate:"required"`
	// Внутренний статус операции
	Status string `json:"status" example:"SUCCESS" validate:"required"`
	// Идентификатор операции на стороне платежной системы
	ExternalID string `json:"external_id,omitempty" example:"ew01r01w0gfw1fw1"`
	// Статус операции на стороне платежной системы
	ExternalStatus string `json:"external_status,omitempty" example:"PENDING"`
	// Время создания операции в формате UNIX Timestamp
	CreatedAt int64 `json:"created_at" example:"1715974447" validate:"required"`
	// Время последнего обновления операции в формате UNIX Timestamp
	UpdatedAt int64 `json:"updated_at" example:"1715974447" validate:"required"`
}

type operationListRequest struct {
	// Идентификатор специалиста поддержки
	UserID int64 `query:"user_id" example:"1" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
}

type operationListResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Массив операций, подходящих под фильтры и условия запроса
	Operations []operation `json:"operations" validate:"required"`
}

// operationList godoc
//
//	@Summary	Получить список операций по заданным фильтрам
//	@Tags		Операции
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		user_id		query		int						true	"Идентификатор специалиста техподдержки"
//	@Param		lang_code	query		string					true	"Код языка, обозначение по RFC 5646"
//	@Success	200			{object}	operationListResponse	"Успешный ответ"
//	@Failure	default		{object}	errorResponse			"Ответ с ошибкой"
//	@Router		/operation/list [get]
func (h *Handler) operationList(c *fiber.Ctx) error {
	ctx := c.Context()

	var req operationListRequest
	if err := c.QueryParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	operations, err := h.operationService.Operations(ctx, req.UserID)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &operationListResponse{
		Success:    true,
		Operations: h.operations(operations),
	}

	return c.JSON(resp)
}

func (h *Handler) operation(item model.Operation) operation {
	return operation{
		ID:             item.ID,
		UserID:         item.UserID,
		Type:           string(item.Type),
		Currency:       item.Currency,
		Amount:         convert.CentsToBase(item.Amount),
		Status:         string(item.Status),
		ExternalID:     item.ExternalID,
		ExternalStatus: string(item.ExternalStatus),
		CreatedAt:      item.CreatedAt.Unix(),
		UpdatedAt:      item.UpdatedAt.Unix(),
	}
}

func (h *Handler) operations(items []model.Operation) []operation {
	var operations []operation
	if itemNum := len(items); itemNum > 0 {
		operations = make([]operation, 0, itemNum)
		for _, item := range items {
			operations = append(operations, h.operation(item))
		}
	}
	return operations
}
