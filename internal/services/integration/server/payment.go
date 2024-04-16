package server

import (
	"context"
	"fmt"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) CreatePayment(ctx context.Context, request *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	integration, ok := s.integrations[request.GetExternalSystem()]
	if !ok || integration == nil {
		return nil, fmt.Errorf("unknown external system: %q", request.GetExternalSystem())
	}

	data := model.CreatePaymentData{
		ReturnURLs:     convert.ReturnURLsFromProto(request.GetReturnUrls()),
		AdditionalData: request.AdditionalData.AsMap(),
		ExternalSystem: request.GetExternalSystem(),
		ExternalMethod: request.GetExternalMethod(),
		Currency:       request.GetCurrency(),
		LangCode:       request.GetLangCode(),
		UserID:         request.GetUserId(),
		Amount:         request.GetAmount(),
		OperationID:    request.GetOperationId(),
	}

	if request.Tool != nil {
		data.Tool = convert.ToolFromProto(request.GetTool())
		data.ToolID = data.Tool.ID
	}

	result, err := integration.CreatePayment(ctx, data)
	if err != nil {
		return nil, err
	}

	response := &pb.CreatePaymentResponse{
		ExternalStatus: convert.OperationExternalStatusToProto(result.ExternalStatus),
	}

	if result.RedirectURL != "" {
		response.RedirectUrl = &result.RedirectURL
	}

	if result.ExternalID != "" {
		response.ExternalId = &result.ExternalID
	}

	if !result.ProcessedAt.IsZero() {
		processedAtUnix := result.ProcessedAt.UTC().Unix()
		response.ProcessedAt = &processedAtUnix
	}

	if result.FailReason != "" {
		response.FailReason = &result.FailReason
	}

	if result.NewAmount > 0 {
		response.NewAmount = &result.NewAmount
	}

	if result.Tool != nil {
		response.Tool = convert.ToolToProto(result.Tool)
	}

	return response, nil
}
