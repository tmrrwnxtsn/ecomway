package server

import (
	"context"
	"fmt"
	"time"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) GetOperationStatus(ctx context.Context, request *pb.GetOperationStatusRequest) (*pb.GetOperationStatusResponse, error) {
	integration, ok := s.integrations[request.GetExternalSystem()]
	if !ok || integration == nil {
		return nil, fmt.Errorf("unknown external system: %q", request.GetExternalSystem())
	}

	data := model.GetOperationStatusData{
		CreatedAt:      time.Unix(request.GetCreatedAt(), 0).UTC(),
		ExternalID:     request.GetExternalId(),
		ExternalSystem: request.GetExternalSystem(),
		ExternalMethod: request.GetExternalMethod(),
		Currency:       request.GetCurrency(),
		OperationType:  convert.OperationTypeFromProto(request.GetOperationType()),
		OperationID:    request.GetOperationId(),
		UserID:         request.GetUserId(),
		Amount:         request.GetAmount(),
	}

	result, err := integration.GetOperationStatus(ctx, data)
	if err != nil {
		return nil, err
	}

	response := &pb.GetOperationStatusResponse{
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

	if result.NewAmount > 0 {
		response.NewAmount = &result.NewAmount
	}

	return response, nil
}
