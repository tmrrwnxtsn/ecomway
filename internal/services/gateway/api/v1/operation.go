package v1

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type operation struct {
	// Идентификатор операции
	ID int64 `json:"id" example:"1" validate:"required"`
	// Тип операции
	Type string `json:"type" example:"payment" validate:"required"`
	// Валюта операции
	Currency string `json:"currency" example:"RUB" validate:"required"`
	// Сумма операции
	Amount float64 `json:"amount" example:"121.01" validate:"required"`
	// Внутренний статус операции
	Status string `json:"status" example:"SUCCESS" validate:"required"`
	// Платежное средство, используемое в операции
	ToolDisplayed string `json:"tool,omitempty" example:"5748********4124"`
	// Время создания операции в формате UNIX Timestamp
	CreatedAt int64 `json:"created_at" example:"1715974447" validate:"required"`
	// Время завершения операции на стороне платежной системы в формате UNIX Timestamp
	ProcessedAt int64 `json:"processed_at,omitempty" example:"1715974447"`
}

type operationListRequest struct {
	// Идентификатор клиента
	UserID int64 `query:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `query:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
	// Идентификатор операции
	ID int64 `query:"id" example:"1"`
	// Тип операции
	Type string `query:"type" example:"payment"`
	// Внутренние статусы операций, перечисленные через запятую
	Statuses string `query:"statuses" example:"SUCCESS,FAILED"`
	// Время создания операции в формате UNIX Timestamp, с которого возвращать результирующие операции
	CreatedAtFrom int64 `query:"created_at_from" example:"1715974447"`
	// Время создания операции в формате UNIX Timestamp, до которого возвращать результирующие операции
	CreatedAtTo int64 `query:"created_at_to" example:"1715974447"`
	// Поле для сортировки результирующего списка (по умолчанию - "id")
	OrderField string `query:"order_field" example:"amount"`
	// Тип сортировки (по умолчанию - "DESC" - по убыванию)
	OrderType string `query:"order_type" example:"ASC"`
}

type operationListResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Сумма всех операций из результирующего массива
	TotalAmount float64 `json:"total_amount" example:"1421.10" validate:"required"`
	// Количество всех операций из результирующего массива
	TotalCount int64 `json:"total_count" example:"15" validate:"required"`
	// Массив операций, подходящих под фильтры и условия запроса
	Operations []operation `json:"operations" validate:"required"`
}

// operationList godoc
//
//	@Summary	Получить список операций по заданным фильтрам
//	@Tags		Операции
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		user_id			query		int						true	"Идентификатор клиента"
//	@Param		session_id		query		string					true	"Идентификатор сессии клиента"
//	@Param		lang_code		query		string					true	"Код языка, обозначение по RFC 5646"
//	@Param		id				query		int						false	"Идентификатор операции"
//	@Param		type			query		string					false	"Тип операции"
//	@Param		statuses		query		string					false	"Внутренние статусы операций, перечисленные через запятую"
//	@Param		created_at_from	query		int						false	"Время создания операции в формате UNIX Timestamp, с которого возвращать результирующие операции"
//	@Param		created_at_to	query		int						false	"Время создания операции в формате UNIX Timestamp, до которого возвращать результирующие операции"
//	@Param		order_field		query		string					false	"Поле для сортировки результирующего списка (по умолчанию - id)"
//	@Param		order_type		query		string					false	"Тип сортировки (по умолчанию - DESC, по убыванию)"
//	@Success	200				{object}	operationListResponse	"Успешный ответ"
//	@Failure	default			{object}	errorResponse			"Ответ с ошибкой"
//	@Router		/operation [get]
func (h *Handler) operationList(c *fiber.Ctx) error {
	ctx := c.Context()

	var req operationListRequest
	if err := c.QueryParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	criteria, err := operationListCriteriaFromRequest(req)
	if err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	operations, err := h.operationService.ReportOperations(ctx, criteria)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	totalAmount, totalCount := h.summaryService.CalculateReportOperationsSummary(operations)

	operations = h.sortingService.SortReportOperations(operations, req.OrderField, req.OrderType)

	respOperations := h.operations(operations)

	resp := &operationListResponse{
		Success:     true,
		TotalAmount: totalAmount,
		TotalCount:  totalCount,
		Operations:  respOperations,
	}

	return c.JSON(resp)
}

func (h *Handler) operation(item model.ReportOperation) operation {
	result := operation{
		ID:            item.ID,
		Type:          string(item.Type),
		Currency:      item.Currency,
		Amount:        convert.CentsToBase(item.Amount),
		Status:        string(item.Status),
		ToolDisplayed: item.ToolDisplayed,
		CreatedAt:     item.CreatedAt.Unix(),
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

func operationListCriteriaFromRequest(req operationListRequest) (model.OperationCriteria, error) {
	criteria := model.OperationCriteria{
		UserID: &req.UserID,
	}

	if req.ID > 0 {
		criteria.ID = &req.ID
	}
	if req.Type != "" {
		switch t := model.OperationType(req.Type); t {
		case model.OperationTypePayment, model.OperationTypePayout:
			criteria.Types = &[]model.OperationType{t}
		default:
			return model.OperationCriteria{}, fmt.Errorf("unresolved operation type: %v", req.Type)
		}
	}
	if req.Statuses != "" {
		criteria.Statuses = &[]model.OperationStatus{}

		statuses := strings.Split(req.Statuses, ",")
		for _, s := range statuses {
			switch status := model.OperationStatus(s); status {
			case model.OperationStatusNew, model.OperationStatusFailed, model.OperationStatusSuccess:
				*criteria.Statuses = append(*criteria.Statuses, status)
			default:
				return model.OperationCriteria{}, fmt.Errorf("unresolved status: %v", s)
			}
		}
	}
	if req.CreatedAtFrom > 0 {
		criteria.CreatedAtFrom = time.Unix(req.CreatedAtFrom, 0).UTC()
	}
	if req.CreatedAtTo > 0 {
		criteria.CreatedAtTo = time.Unix(req.CreatedAtTo, 0).UTC()
	}

	return criteria, nil
}
