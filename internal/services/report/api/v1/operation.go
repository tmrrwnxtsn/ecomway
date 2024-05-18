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
	// Платежное средство, используемое в операции
	ToolDisplayed string `json:"tool,omitempty" example:"5748********4124"`
	// Причина отклонения операции
	FailReason string `json:"fail_reason,omitempty" example:"Technical error"`
	// Время создания операции в формате UNIX Timestamp
	CreatedAt int64 `json:"created_at" example:"1715974447" validate:"required"`
	// Время последнего обновления операции в формате UNIX Timestamp
	UpdatedAt int64 `json:"updated_at" example:"1715974447" validate:"required"`
	// Время завершения операции на стороне платежной системы в формате UNIX Timestamp
	ProcessedAt int64 `json:"processed_at,omitempty" example:"1715974447"`
}

type operationListRequest struct {
	// Идентификатор специалиста поддержки
	UserID int64 `query:"user_id" example:"1" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
	// Поле для сортировки результирующего списка (по умолчанию - "id")
	OrderField string `query:"order_field" example:"amount"`
	// Тип сортировки (по умолчанию - "DESC" - по убыванию)
	OrderType string `query:"order_type" example:"ASC"`
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
//	@Param		order_field	query		string					false	"Поле для сортировки результирующего списка (по умолчанию - id)"
//	@Param		order_type	query		string					false	"Тип сортировки (по умолчанию - DESC, по убыванию)"
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

	operations, err := h.operationService.ReportOperations(ctx, req.UserID)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	operations = h.sortingService.SortReportOperations(operations, req.OrderField, req.OrderType)

	resp := &operationListResponse{
		Success:    true,
		Operations: h.operations(operations),
	}

	return c.JSON(resp)
}

func (h *Handler) operation(item model.ReportOperation) operation {
	result := operation{
		ID:             item.ID,
		UserID:         item.UserID,
		Type:           string(item.Type),
		Currency:       item.Currency,
		Amount:         convert.CentsToBase(item.Amount),
		Status:         string(item.Status),
		ExternalID:     item.ExternalID,
		ExternalStatus: string(item.ExternalStatus),
		ToolDisplayed:  item.ToolDisplayed,
		FailReason:     item.FailReason,
		CreatedAt:      item.CreatedAt.Unix(),
		UpdatedAt:      item.UpdatedAt.Unix(),
	}

	if !item.ProcessedAt.IsZero() {
		result.ProcessedAt = item.ProcessedAt.Unix()
	}

	return result
}

func (h *Handler) operations(items []model.ReportOperation) []operation {
	var operations []operation
	if itemNum := len(items); itemNum > 0 {
		operations = make([]operation, 0, itemNum)
		for _, item := range items {
			operations = append(operations, h.operation(item))
		}
	}
	return operations
}
