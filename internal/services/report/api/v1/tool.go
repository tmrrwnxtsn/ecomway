package v1

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/translate"
)

const (
	toolTypeBankCard = "card"
	toolTypeWallet   = "wallet"
)

type toolDetails struct {
	// Тип банковской карты
	CardType string `json:"card_type,omitempty" example:"Visa"`
	// Владелец банковской карты
	CardHolder string `json:"card_holder,omitempty" example:"Ivanov Ivan"`
	// Срок действия банковской карты (месяц, MM)
	ExpiryMonth int64 `json:"expiry_month,omitempty" example:"10"`
	// Срок действия банковской карты (год, YYYY)
	ExpiryYear int64 `json:"expiry_year,omitempty" example:"2023"`
	// Название банка, выпустившего банковскую карту
	BankName string `json:"bank_name,omitempty" example:"Sberbank"`
	// Номер электронного кошелька
	WalletNumber string `json:"wallet_number,omitempty" example:"410011758831136"`
}

type tool struct {
	// Идентификатор платежного средства
	ID string `json:"id" example:"2dc32aa0-000f-5000-8000-16d7bc6cd09f" validate:"required"`
	// Внутренний код платежного метода платежной системы, к которой относится платежное средство
	ExternalMethod string `json:"external_method" example:"yookassa_bank_card" validate:"required"`
	// Тип платежного средства:
	// * Банковская карта - "card"
	// * Электронный кошелек - "wallet"
	Type string `json:"type" example:"card" validate:"required"`
	// Статус платежного средства:
	// * Доступен клиенту - "ACTIVE"
	// * Удален клиентом - "REMOVED_BY_CLIENT"
	// * Ожидает восстановления - "PENDING_RECOVERY"
	// * Заблокирован техподдержкой - "REMOVED_BY_ADMINISTRATOR"
	Status string `json:"status" example:"ACTIVE" validate:"required"`
	// Название платежного средства
	Name string `json:"name" example:"Карта брата" validate:"required"`
	// Значение платежного средства:
	// * Маскированная банковская карта
	// * Номер электронного кошелька
	Caption string `json:"caption" example:"444444******4444" validate:"required"`
	// Время создания платежного средства в формате UNIX Timestamp
	CreatedAt int64 `json:"created_at" example:"1715974447" validate:"required"`
	// Время последнего обновления платежного средства в формате UNIX Timestamp
	UpdatedAt int64 `json:"updated_at" example:"1715974447" validate:"required"`
	// Дополнительная информация о платежном средстве
	Details *toolDetails `json:"details,omitempty"`
}

type toolListRequest struct {
	// Идентификатор специалиста поддержки
	UserID int64 `query:"user_id" example:"1" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
	// Идентификатор клиента
	ClientID int64 `query:"client_id" example:"1" validate:"required"`
}

type toolListResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Массив платежных средств клиента
	Tools []tool `json:"tools" validate:"required"`
}

// toolList godoc
//
//	@Summary	Получить список платежных средств клиента
//	@Tags		Платежные средства
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		user_id		query		int					true	"Идентификатор специалиста техподдержки"
//	@Param		lang_code	query		string				true	"Код языка, обозначение по RFC 5646"
//	@Param		client_id	query		int					true	"Идентификатор клиента"
//	@Success	200			{object}	toolListResponse	"Успешный ответ"
//	@Failure	default		{object}	errorResponse		"Ответ с ошибкой"
//	@Router		/tool [get]
func (h *Handler) toolList(c *fiber.Ctx) error {
	ctx := c.Context()

	var req toolListRequest
	if err := c.QueryParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	tools, err := h.toolService.AllTools(ctx, req.ClientID)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &toolListResponse{
		Success: true,
		Tools:   h.tools(tools),
	}

	return c.JSON(resp)
}

type toolRecoverRequest struct {
	// Идентификатор специалиста поддержки
	UserID int64 `json:"user_id" example:"1" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `json:"lang_code" example:"en" validate:"required"`
	// Идентификатор платежного средства
	ID string `json:"id" example:"2dc32aa0-000f-5000-8000-16d7bc6cd09f" validate:"required"`
	// Внутренний код платежного метода платежной системы, к которой относится платежное средство
	ExternalMethod string `json:"external_method" example:"yookassa_bank_card" validate:"required"`
	// Идентификатор клиента
	ClientID int64 `json:"client_id" example:"1" validate:"required"`
}

type toolRecoverResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Сообщение, которое необходимо показать специалисту техподдержки
	Message string `json:"message" example:"Платежное средство удалено." validate:"required"`
}

// toolRecover godoc
//
//	@Summary	Установить платежное средство готовым к восстановлению
//	@Tags		Платежные средства
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		input	body		toolRecoverRequest	true	"Тело запроса"
//	@Success	200		{object}	toolRecoverResponse	"Успешный ответ"
//	@Failure	default	{object}	errorResponse		"Ответ с ошибкой"
//	@Router		/tool/recover [put]
func (h *Handler) toolRecover(c *fiber.Ctx) error {
	ctx := c.Context()

	var req toolRecoverRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.toolService.RecoverTool(ctx, req.ID, req.ClientID, req.ExternalMethod); err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal {
				switch perr.Code {
				case perror.CodeObjectNotFound:
					return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
				case perror.CodeUnresolvedStatusConflict:
					return h.unresolvedObjectStatusErrorResponse(c, req.LangCode, perr)
				}
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &toolRecoverResponse{
		Success: true,
		Message: h.translator.Translate(req.LangCode, translate.KeyToolRecovered),
	}

	return c.JSON(resp)
}

func (h *Handler) tool(item *model.Tool) tool {
	t := tool{
		ID:             item.ID,
		ExternalMethod: item.ExternalMethod,
		Status:         string(item.Status),
		Name:           item.Name,
		Caption:        item.Displayed,
		CreatedAt:      item.CreatedAt.UTC().Unix(),
		UpdatedAt:      item.UpdatedAt.UTC().Unix(),
	}

	switch item.Type {
	case model.ToolTypeBankCard:
		t.Type = toolTypeBankCard

		cardType, _ := item.Details["card_type"].(string)
		cardHolder, _ := item.Details["card_holder"].(string)
		bankName, _ := item.Details["bank_name"].(string)
		expiryMonthStr, _ := item.Details["expiry_month"].(string)
		expiryYearStr, _ := item.Details["expiry_year"].(string)

		if cardType+cardHolder+bankName+expiryMonthStr+expiryYearStr == "" {
			break
		}

		expiryMonth, _ := strconv.ParseInt(expiryMonthStr, 10, 64)
		expiryYear, _ := strconv.ParseInt(expiryYearStr, 10, 64)

		t.Details = &toolDetails{
			CardType:    cardType,
			CardHolder:  cardHolder,
			ExpiryMonth: expiryMonth,
			ExpiryYear:  expiryYear,
			BankName:    bankName,
		}
	case model.ToolTypeWallet:
		t.Type = toolTypeWallet

		walletNumber, _ := item.Details["number"].(string)
		if walletNumber == "" {
			break
		}

		t.Details = &toolDetails{
			WalletNumber: walletNumber,
		}
	}

	return t
}

func (h *Handler) tools(items []*model.Tool) []tool {
	var tools []tool
	if itemNum := len(items); itemNum > 0 {
		tools = make([]tool, 0, itemNum)
		for _, item := range items {
			tools = append(tools, h.tool(item))
		}
	}
	return tools
}
