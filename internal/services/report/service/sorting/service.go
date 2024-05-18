package sorting

import (
	"sort"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

const (
	orderFieldID        = "id"
	orderFieldUserID    = "client_id"
	orderFieldType      = "type"
	orderFieldAmount    = "amount"
	orderFieldStatus    = "status"
	orderFieldCreatedAt = "created_at"
	orderFieldUpdatedAt = "updated_at"
)

const (
	orderTypeASC  = "ASC"
	orderTypeDESC = "DESC"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SortReportOperations(items []model.ReportOperation, orderField, orderType string) []model.ReportOperation {
	if orderField == "" {
		orderField = orderFieldID
	}
	if orderType == "" {
		orderType = orderTypeDESC
	}

	sort.Slice(items, func(i, j int) bool {
		switch orderField {
		case orderFieldUserID:
			if orderType == orderTypeASC {
				return items[i].UserID < items[j].UserID
			}
			return items[i].UserID > items[j].UserID
		case orderFieldType:
			if orderType == orderTypeASC {
				return items[i].Type < items[j].Type
			}
			return items[i].Type > items[j].Type
		case orderFieldAmount:
			if orderType == orderTypeASC {
				return items[i].Amount < items[j].Amount
			}
			return items[i].Amount > items[j].Amount
		case orderFieldStatus:
			if orderType == orderTypeASC {
				return items[i].Status < items[j].Status
			}
			return items[i].Status > items[j].Status
		case orderFieldCreatedAt:
			if orderType == orderTypeASC {
				return items[i].CreatedAt.Before(items[j].CreatedAt)
			}
			return items[i].CreatedAt.After(items[j].CreatedAt)
		case orderFieldUpdatedAt:
			if orderType == orderTypeASC {
				return items[i].UpdatedAt.Before(items[j].UpdatedAt)
			}
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		default:
			if orderType == orderTypeASC {
				return items[i].ID < items[j].ID
			}
			return items[i].ID > items[j].ID
		}
	})

	return items
}
