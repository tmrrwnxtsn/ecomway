package engine

import (
	"context"
	"time"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) ReportOperations(ctx context.Context, userID int64) ([]model.ReportOperation, error) {
	request := &pb.ReportOperationsRequest{
		UserId: userID,
	}

	response, err := c.client.ReportOperations(ctx, request)
	if err != nil {
		return nil, err
	}

	operations := make([]model.ReportOperation, 0, len(response.GetOperations()))
	for _, pbOp := range response.GetOperations() {
		operations = append(operations, reportOperationFromProto(pbOp))
	}
	return operations, nil
}

func reportOperationFromProto(op *pb.ReportOperation) model.ReportOperation {
	result := model.ReportOperation{
		ID:             op.GetId(),
		UserID:         op.GetUserId(),
		Type:           convert.OperationTypeFromProto(op.GetType()),
		Currency:       op.GetCurrency(),
		Amount:         op.GetAmount(),
		Status:         convert.OperationStatusFromProto(op.GetStatus()),
		ExternalID:     op.GetExternalId(),
		ExternalSystem: op.GetExternalSystem(),
		ExternalMethod: op.GetExternalMethod(),
		ExternalStatus: convert.OperationExternalStatusFromProto(op.GetExternalStatus()),
		ToolDisplayed:  op.GetToolDisplayed(),
		FailReason:     op.GetFailReason(),
		CreatedAt:      time.Unix(op.GetCreatedAt(), 0).UTC(),
		UpdatedAt:      time.Unix(op.GetUpdatedAt(), 0).UTC(),
		ProcessedAt:    time.Time{},
	}

	if op.ProcessedAt != nil {
		result.ProcessedAt = time.Unix(op.GetProcessedAt(), 0).UTC()
	}

	return result
}
