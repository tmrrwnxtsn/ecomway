package server

import (
	"context"
	"fmt"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) CreatePayout(ctx context.Context, request *pb.CreatePayoutRequest) (*pb.CreatePayoutResponse, error) {
	integration, ok := s.integrations[request.GetExternalSystem()]
	if !ok || integration == nil {
		return nil, fmt.Errorf("unknown external system: %q", request.GetExternalSystem())
	}

	tool := convert.ToolFromProto(request.GetTool())

	data := model.CreatePayoutData{
		Tool:           tool,
		AdditionalData: request.AdditionalData.AsMap(),
		ExternalSystem: request.GetExternalSystem(),
		ExternalMethod: request.GetExternalMethod(),
		Currency:       request.GetCurrency(),
		LangCode:       request.GetLangCode(),
		UserID:         request.GetUserId(),
		ToolID:         tool.ID,
		Amount:         request.GetAmount(),
		OperationID:    request.GetOperationId(),
	}

	result, err := integration.CreatePayout(ctx, data)
	if err != nil {
		return nil, err
	}

	response := &pb.CreatePayoutResponse{
		ExternalStatus: convert.OperationExternalStatusToProto(result.ExternalStatus),
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

	return response, nil
}
