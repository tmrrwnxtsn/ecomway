package server

import (
	"context"
	"fmt"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) CreatePayment(ctx context.Context, request *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	integration, ok := s.integrations[request.GetExternalSystem()]
	if !ok || integration == nil {
		return nil, fmt.Errorf("unknown external system: %q", request.GetExternalSystem())
	}

	data := model.CreatePaymentData{
		AdditionalData: request.AdditionalData.AsMap(),
		ExternalSystem: request.GetExternalSystem(),
		ExternalMethod: request.GetExternalMethod(),
		Currency:       request.GetCurrency(),
		LangCode:       request.GetLangCode(),
		UserID:         request.GetUserId(),
		Amount:         request.GetAmount(),
		OperationID:    request.GetOperationId(),
	}

	result, err := integration.CreatePayment(ctx, data)
	if err != nil {
		return nil, err
	}

	response := &pb.CreatePaymentResponse{
		RedirectUrl: result.RedirectURL,
	}

	if result.ExternalID != "" {
		response.ExternalId = &result.ExternalID
	}

	return response, nil
}