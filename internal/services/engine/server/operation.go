package server

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) ReportOperations(ctx context.Context, request *pb.ReportOperationsRequest) (*pb.ReportOperationsResponse, error) {
	criteria := criteriaFromReportOperationsRequest(request)

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

func (s *Server) GetOperationExternalStatus(ctx context.Context, request *pb.GetOperationExternalStatusRequest) (*pb.GetOperationExternalStatusResponse, error) {
	criteria := model.OperationCriteria{
		ID: &request.OperationId,
	}

	operation, err := s.operationService.GetOne(ctx, criteria)
	if err != nil {
		return nil, err
	}

	statusData := model.GetOperationStatusData{
		CreatedAt:      operation.CreatedAt,
		ExternalID:     operation.ExternalID,
		ExternalSystem: operation.ExternalSystem,
		ExternalMethod: operation.ExternalMethod,
		Currency:       operation.Currency,
		OperationType:  operation.Type,
		OperationID:    operation.ID,
		UserID:         operation.UserID,
		Amount:         operation.Amount,
	}

	result, err := s.integrationClient.GetOperationStatus(ctx, statusData)
	if err != nil {
		return nil, err
	}

	return &pb.GetOperationExternalStatusResponse{
		ExternalStatus: convert.OperationExternalStatusToProto(result.ExternalStatus),
	}, nil
}

func (s *Server) ChangeOperationStatus(ctx context.Context, request *pb.ChangeOperationStatusRequest) (*emptypb.Empty, error) {
	id := request.GetId()

	criteria := model.OperationCriteria{
		ID: &id,
	}

	operation, err := s.operationService.GetOne(ctx, criteria)
	if err != nil {
		return nil, err
	}

	newStatus := convert.OperationStatusFromProto(request.GetNewStatus())
	newExternalStatus := convert.OperationExternalStatusFromProto(request.GetNewExternalStatus())

	result, err := s.operationService.ChangeStatus(ctx, id, newStatus, newExternalStatus)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return &emptypb.Empty{}, nil
	}

	switch result {
	case model.OperationChangeStatusResultFailPayment:
		err = s.paymentService.Fail(ctx, model.FailPaymentData{
			ExternalID:     operation.ExternalID,
			ExternalStatus: newExternalStatus,
			FailReason:     model.OperationFailReasonManual,
			OperationID:    operation.ID,
		})
	case model.OperationChangeStatusResultFailPayout:
		err = s.payoutService.Fail(ctx, model.FailPayoutData{
			ExternalID:     operation.ExternalID,
			ExternalStatus: newExternalStatus,
			FailReason:     model.OperationFailReasonManual,
			OperationID:    operation.ID,
		})
	case model.OperationChangeStatusResultSuccessPayment:
		err = s.paymentService.Success(ctx, model.SuccessPaymentData{
			ProcessedAt:    time.Now(),
			ExternalID:     operation.ExternalID,
			ExternalStatus: operation.ExternalStatus,
			OperationID:    operation.ID,
			NewAmount:      operation.Amount,
		})
	case model.OperationChangeStatusResultSuccessPayout:
		tool, err := s.toolService.GetOne(ctx, operation.ToolID, operation.UserID, operation.ExternalMethod)
		if err != nil {
			return nil, fmt.Errorf("get tool from db: %w", err)
		}

		err = s.payoutService.Success(ctx, model.SuccessPayoutData{
			ProcessedAt:    time.Now(),
			ExternalID:     operation.ExternalID,
			ExternalStatus: newExternalStatus,
			OperationID:    operation.ID,
			Tool:           tool,
		})
	}
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func criteriaFromReportOperationsRequest(request *pb.ReportOperationsRequest) model.OperationCriteria {
	criteria := model.OperationCriteria{
		ID:         request.Id,
		UserID:     request.UserId,
		ExternalID: request.ExternalId,
	}
	if len(request.Types) > 0 {
		types := make([]model.OperationType, 0, len(request.Types))
		for _, pbType := range request.Types {
			types = append(types, convert.OperationTypeFromProto(pbType))
		}
		criteria.Types = &types
	}
	if len(request.Statuses) > 0 {
		statuses := make([]model.OperationStatus, 0, len(request.Statuses))
		for _, pbStatus := range request.Statuses {
			statuses = append(statuses, convert.OperationStatusFromProto(pbStatus))
		}
		criteria.Statuses = &statuses
	}
	if request.CreatedAtFrom != nil {
		criteria.CreatedAtFrom = time.Unix(request.GetCreatedAtFrom(), 0).UTC()
	}
	if request.CreatedAtTo != nil {
		criteria.CreatedAtTo = time.Unix(request.GetCreatedAtTo(), 0).UTC()
	}
	return criteria
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
