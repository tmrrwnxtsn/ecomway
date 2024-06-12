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
	// Название платежного средства
	Name string `json:"name" example:"Карта брата" validate:"required"`
	// Тип платежного средства:
	// * Банковская карта - "card"
	// * Электронный кошелек - "wallet"
	Type string `json:"type" example:"card" validate:"required"`
	// Значение платежного средства, например:
	// * Маскированная банковская карта
	// * Номер электронного кошелька
	// * Адрес электронной почты
	// * и т.д.
	Caption string `json:"caption" example:"444444******4444" validate:"required"`
	// Дополнительная информация о платежном средстве
	Details *toolDetails `json:"details,omitempty"`
}

type toolListRequest struct {
	// Идентификатор клиента
	UserID int64 `query:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `query:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
}

type toolListResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Массив сохраненных платежных средств клиента
	Tools []tool `json:"tools" validate:"required"`
}

// toolList godoc
//
//	@Summary	Получить список сохраненных платежных средств клиента
//	@Tags		Платежные средства
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		user_id		query		int					true	"Идентификатор клиента"
//	@Param		session_id	query		string				true	"Идентификатор сессии клиента"
//	@Param		lang_code	query		string				true	"Код языка, обозначение по RFC 5646"
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

	tools, err := h.toolService.AvailableTools(ctx, req.UserID)
	if err != nil {
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &toolListResponse{
		Success: true,
		Tools:   h.tools(tools),
	}

	return c.JSON(resp)
}

type toolEditRequest struct {
	// Идентификатор платежного средства
	ID string `json:"id" example:"2dc32aa0-000f-5000-8000-16d7bc6cd09f" validate:"required"`
	// Идентификатор клиента
	UserID int64 `json:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `json:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
	// Внутренний код платежного метода платежной системы, к которой относится платежное средство
	ExternalMethod string `json:"external_method" example:"yookassa_bank_card" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `json:"lang_code" example:"en" validate:"required"`
	// Новое название платежного инструмента, выбранное клиентом
	Name string `json:"name" example:"Карта брата" validate:"required"`
}

type toolEditResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Информация об измененном платежном средстве
	Tool tool `json:"tool" validate:"required"`
}

// toolEdit godoc
//
//	@Summary	Изменить информацию о платежном средстве
//	@Tags		Платежные средства
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		input	body		toolEditRequest		true	"Тело запроса"
//	@Success	200		{object}	toolEditResponse	"Успешный ответ"
//	@Failure	default	{object}	errorResponse		"Ответ с ошибкой"
//	@Router		/tool/edit [put]
func (h *Handler) toolEdit(c *fiber.Ctx) error {
	ctx := c.Context()

	var req toolEditRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	edited, err := h.toolService.EditTool(ctx, req.ID, req.UserID, req.ExternalMethod, req.Name)
	if err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal && perr.Code == perror.CodeObjectNotFound {
				return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &toolEditResponse{
		Success: true,
		Tool:    h.tool(edited),
	}

	return c.JSON(resp)
}

type toolRemoveRequest struct {
	// Идентификатор платежного средства
	ID string `json:"id" example:"2dc32aa0-000f-5000-8000-16d7bc6cd09f" validate:"required"`
	// Идентификатор клиента
	UserID int64 `json:"user_id" example:"1" validate:"required"`
	// Идентификатор сессии клиента
	SessionID string `json:"session_id" example:"LRXZmXPGusPCfys48LadjFew" validate:"required"`
	// Внутренний код платежного метода платежной системы, к которой относится платежное средство
	ExternalMethod string `json:"external_method" example:"yookassa_bank_card" validate:"required"`
	// Код языка, обозначение по RFC 5646
	LangCode string `json:"lang_code" example:"en" validate:"required"`
}

type toolRemoveResponse struct {
	// Результат обработки запроса (всегда true)
	Success bool `json:"success" example:"true" validate:"required"`
	// Сообщение, которое необходимо показать клиенту
	Message string `json:"message" example:"Платежное средство удалено." validate:"required"`
}

// toolRemove godoc
//
//	@Summary	Удалить платежное средство
//	@Tags		Платежные средства
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		input	body		toolRemoveRequest	true	"Тело запроса"
//	@Success	200		{object}	toolRemoveResponse	"Успешный ответ"
//	@Failure	default	{object}	errorResponse		"Ответ с ошибкой"
//	@Router		/tool/remove [delete]
func (h *Handler) toolRemove(c *fiber.Ctx) error {
	ctx := c.Context()

	var req toolRemoveRequest
	if err := c.BodyParser(&req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return h.requestValidationErrorResponse(c, req.LangCode, err)
	}

	if err := h.toolService.RemoveTool(ctx, req.ID, req.UserID, req.ExternalMethod); err != nil {
		var perr *perror.Error
		if errors.As(err, &perr) {
			if perr.Group == perror.GroupInternal {
				switch perr.Code {
				case perror.CodeObjectNotFound:
					return h.objectNotFoundErrorResponse(c, req.LangCode, perr)
				}
			}
		}
		return h.internalErrorResponse(c, req.LangCode, err)
	}

	resp := &toolRemoveResponse{
		Success: true,
		Message: h.translator.Translate(req.LangCode, translate.KeyToolRemoved),
	}

	return c.JSON(resp)
}

func (h *Handler) tool(item *model.Tool) tool {
	t := tool{
		ID:             item.ID,
		ExternalMethod: item.ExternalMethod,
		Caption:        item.Displayed,
		Name:           item.Name,
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
