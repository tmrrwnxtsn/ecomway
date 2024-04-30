package v1

import (
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type limits struct {
	// Код валюты лимита в соответствии со стандартом ISO 4217
	Currency string `json:"currency" example:"RUB" validate:"required"`
	// Минимальное значение суммы (в дробных единицах)
	MinAmount float64 `json:"min_amount" example:"100" validate:"required"`
	// Максимальное значение суммы (в дробных единицах)
	MaxAmount float64 `json:"max_amount" example:"60000" validate:"required"`
}

type commission struct {
	// Тип комиссии:
	// * "percent" - комиссия в процентах
	// * "fixed" - фиксированная комиссия
	// * "combined" - комбинированная комиссия, из процентов и фиксированной суммы
	// * "text" - текстовая коммиссия, например, "Взимается провайдером"
	Type string `json:"type" example:"combined" validate:"required"`
	// Значение комиссии (передается, если тип комиссии, "type", равен "percent" или "combined")
	Percent *float64 `json:"percent,omitempty" example:"11.99" validate:"optional"`
	// Значение комиссии (передается, если тип комиссии, "type", равен "fixed" или "combined")
	Absolute *float64 `json:"absolute,omitempty" example:"10" validate:"optional"`
	// Код валюты комиссии (передается, если тип комиссии, "type", равен "fixed" или "combined")
	Currency *string `json:"currency,omitempty" example:"RUB" validate:"optional"`
	// Текстовая репрезентация комиссии (передается, если тип комиссии, "type", равен "text")
	Caption *string `json:"caption,omitempty" example:"Комиссия взимается провайдером" validate:"optional"`
}

type method struct {
	// Идентификатор ПС из внутреннего справочника
	ID string `json:"id" example:"CARD" validate:"required"`
	// Название платежной системы
	Name string `json:"name" example:"Банковская карта" validate:"required"`
	// Внутренний код платежной системы, к которой направляется целевой запрос
	ExternalSystem string `json:"external_system" example:"yookassa" validate:"required"`
	// Внутренний код платежного метода платежной системы, к которой направляется целевой запрос
	ExternalMethod string `json:"external_method" example:"yookassa_bank_card" validate:"required"`
	// Флаг о том, что платежная система добавлена в избранное
	IsFavorite bool `json:"favorite" example:"true" validate:"required"`
	// Массив объектов, содержащих данные о лимитах
	Limits []limits `json:"limits" validate:"required"`
	// Объект, содержащий данные о комиссии
	Commission commission `json:"commission" validate:"required"`
	// Массив объектов, содержащих данные о сохраненных платежных средствах
	Tools []tool `json:"tools,omitempty"`
}

func (h *Handler) methods(items []model.Method, toolsGrouped map[string][]*model.Tool, langCode string) []method {
	result := make([]method, 0, len(items))
	for _, item := range items {
		tools := toolsGrouped[item.ExternalMethod]

		result = append(result, h.method(item, tools, langCode))
	}
	return result
}

func (h *Handler) method(item model.Method, tools []*model.Tool, langCode string) method {
	return method{
		ID:             item.ID,
		Name:           item.DisplayedName[langCode],
		ExternalSystem: item.ExternalSystem,
		ExternalMethod: item.ExternalMethod,
		IsFavorite:     false,
		Limits:         h.limits(item.Limits),
		Commission:     h.commission(item.Commission, langCode),
		Tools:          h.tools(tools),
	}
}

func (h *Handler) limits(items map[string]model.Limits) []limits {
	result := make([]limits, 0, len(items))
	for currency, l := range items {
		result = append(result, limits{
			Currency:  currency,
			MinAmount: convert.CentsToBase(l.MinAmount),
			MaxAmount: convert.CentsToBase(l.MaxAmount),
		})
	}
	return result
}

func (h *Handler) commission(item model.Commission, langCode string) commission {
	result := commission{
		Type:     string(item.Type),
		Percent:  item.Percent,
		Absolute: item.Absolute,
	}

	if item.Currency != "" {
		result.Currency = &item.Currency
	}

	caption, ok := item.Message[langCode]
	if ok {
		result.Caption = &caption
	}

	return result
}
