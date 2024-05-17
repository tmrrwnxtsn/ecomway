package server

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) ReportOperations(ctx context.Context, _ *pb.ReportOperationsRequest) (*pb.ReportOperationsResponse, error) {
	criteria := model.OperationCriteria{}

	operations, err := s.operationService.AllForReport(ctx, criteria)
	if err != nil {
		return nil, err
	}

	pbOperations := make([]*pb.ReportOperation, 0, len(operations))
	for _, op := range operations {
		pbOperations = append(pbOperations, reportOperationToProto(op))
	}

	return &pb.ReportOperationsResponse{
		Operations: pbOperations,
	}, nil
}

func reportOperationToProto(op model.ReportOperation) *pb.ReportOperation {
	result := &pb.ReportOperation{
		Id:             op.ID,
		UserId:         op.UserID,
		Type:           convert.OperationTypeToProto(op.Type),
		Currency:       op.Currency,
		Amount:         op.Amount,
		Status:         convert.OperationStatusToProto(op.Status),
		ExternalSystem: op.ExternalSystem,
		ExternalMethod: op.ExternalMethod,
		CreatedAt:      op.CreatedAt.UTC().Unix(),
		UpdatedAt:      op.UpdatedAt.UTC().Unix(),
	}

	if op.ExternalID != "" {
		result.ExternalId = &op.ExternalID
	}

	if op.ExternalStatus != "" {
		pbExternalStatus := convert.OperationExternalStatusToProto(op.ExternalStatus)
		result.ExternalStatus = &pbExternalStatus
	}

	if op.ToolDisplayed != "" {
		result.ToolDisplayed = &op.ToolDisplayed
	}

	if op.FailReason != "" {
		result.FailReason = &op.FailReason
	}

	if !op.ProcessedAt.IsZero() {
		processedAtUnix := op.ProcessedAt.UTC().Unix()
		result.ProcessedAt = &processedAtUnix
	}

	return result
}
