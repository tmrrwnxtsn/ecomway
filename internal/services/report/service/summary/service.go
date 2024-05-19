package summary

import (
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CalculateReportOperationsSummary(items []model.ReportOperation) (totalAmount float64, totalCount int64) {
	for i := 0; i < len(items); i++ {
		totalAmount += convert.CentsToBase(items[i].Amount)
		totalCount++
	}
	return totalAmount, totalCount
}
