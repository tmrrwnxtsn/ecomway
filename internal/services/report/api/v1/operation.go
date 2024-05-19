package v1

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/translate"
)

type operation struct {
	// Идентификатор операции
	ID int64 `json:"id" example:"1" validate:"required"`
	// Идентификатор клиента
	ClientID int64 `json:"client_id" example:"1" validate:"required"`
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
	// Идентификатор операции
	ID int64 `query:"id" example:"1"`
	// Идентификатор операции на стороне платежной системы
	ExternalID string `query:"external_id" example:"ew01r01w0gfw1fw1"`
	// Идентификатор клиента
	ClientID int64 `query:"client_id" example:"1"`
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
	// Массив операций, подходящих под фильтры и условия запроса
	Operations []operation `json:"operations" validate:"required"`
}

// operationList godoc
//
//	@Summary	Получить список операций по заданным фильтрам
//	@Tags		Операции
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		user_id			query		int						true	"Идентификатор специалиста техподдержки"
//	@Param		lang_code		query		string					true	"Код языка, обозначение по RFC 5646"
//	@Param		id				query		int						false	"Идентификатор операции"
//	@Param		external_id		query		string					false	"Идентификатор операции на стороне платежной системы"
//	@Param		client_id		query		int						false	"Идентификатор клиента"
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

	operations = h.sortingService.SortReportOperations(operations, req.OrderField, req.OrderType)

	resp := &operationListResponse{
		Success:    true,
		Operations: h.operations(operations),
	}

	return c.JSON(resp)
}

func operationListCriteriaFromRequest(req operationListRequest) (model.OperationCriteria, error) {
	var criteria model.OperationCriteria

	if req.ID > 0 {
		criteria.ID = &req.ID
	}
	if req.ClientID > 0 {
		criteria.UserID = &req.ClientID
	}
	if req.ExternalID != "" {
		criteria.ExternalID = &req.ExternalID
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

type operationExternalStatusRequest struct {
	// Идентификатор специалиста поддержки
	UserID int64 `query:"user_id" example:"1" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
}

type operationExternalStatusResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Статус операции на стороне платежной системы
	ExternalStatus string `json:"external_status" example:"PENDING" validate:"required"`
	// Информативное сообщение, описывающее статус транзакции на стороне платежной системы
	Message string `json:"message" example:"Транзакция на стороне ПС еще не имеет конечный статус." validate:"required"`
}

// operationExternalStatus godoc
//
//	@Summary	Запросить статус операции на стороне платежной системы
//	@Tags		Операции
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id			path		int								true	"Идентификатор операции"
//	@Param		user_id		query		int								true	"Идентификатор специалиста техподдержки"
//	@Param		lang_code	query		string							true	"Код языка, обозначение по RFC 5646"
//	@Success	200			{object}	operationExternalStatusResponse	"Успешный ответ"
//	@Failure	default		{object}	errorResponse					"Ответ с ошибкой"
//	@Router		/operation/{id}/external-status [get]
func (h *Handler) operationExternalStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	var req operationExternalStatusRequest
	if err := c.QueryParser(&req); err != nil {
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

	externalOperationStatus, err := h.operationService.GetExternalOperationStatus(ctx, opID)
	if err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal && perr.Code == perror.CodeObjectNotFound {
				return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &operationExternalStatusResponse{
		Success:        true,
		ExternalStatus: string(externalOperationStatus),
		Message:        h.operationExternalStatusMessage(externalOperationStatus, req.LangCode),
	}

	return c.JSON(resp)
}

func (h *Handler) operationExternalStatusMessage(externalStatus model.OperationExternalStatus, langCode string) string {
	switch externalStatus {
	case model.OperationExternalStatusPending:
		return h.translator.Translate(langCode, translate.KeyExternalStatusPending)
	case model.OperationExternalStatusSuccess:
		return h.translator.Translate(langCode, translate.KeyExternalStatusSuccess)
	case model.OperationExternalStatusFailed:
		return h.translator.Translate(langCode, translate.KeyExternalStatusFailed)
	default:
		return h.translator.Translate(langCode, translate.KeyExternalStatusUnknown)
	}
}

func (h *Handler) operation(item model.ReportOperation) operation {
	result := operation{
		ID:             item.ID,
		ClientID:       item.UserID,
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
