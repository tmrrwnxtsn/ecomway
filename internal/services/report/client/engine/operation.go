package engine

import (
	"context"
	"time"

	pbEngine "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) ReportOperations(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error) {
	request := reportOperationRequestFromCriteria(criteria)

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

func (c *Client) GetExternalOperationStatus(ctx context.Context, id int64) (model.OperationExternalStatus, error) {
	request := &pbEngine.GetOperationExternalStatusRequest{
		OperationId: id,
	}

	response, err := c.client.GetOperationExternalStatus(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return "", perr
		}
		return "", err
	}

	externalStatus := convert.OperationExternalStatusFromProto(response.GetExternalStatus())
	return externalStatus, nil
}

func reportOperationRequestFromCriteria(criteria model.OperationCriteria) *pbEngine.ReportOperationsRequest {
	request := &pbEngine.ReportOperationsRequest{
		Id:         criteria.ID,
		UserId:     criteria.UserID,
		ExternalId: criteria.ExternalID,
	}
	if criteria.Types != nil {
		pbTypes := make([]pb.OperationType, 0, len(*criteria.Types))
		for _, operationType := range *criteria.Types {
			pbTypes = append(pbTypes, convert.OperationTypeToProto(operationType))
		}
		request.Types = pbTypes
	}
	if criteria.Statuses != nil {
		pbStatuses := make([]pb.OperationStatus, 0, len(*criteria.Statuses))
		for _, operationStatus := range *criteria.Statuses {
			pbStatuses = append(pbStatuses, convert.OperationStatusToProto(operationStatus))
		}
		request.Statuses = pbStatuses
	}
	if !criteria.CreatedAtFrom.IsZero() {
		createdAtFromUnix := criteria.CreatedAtFrom.UTC().Unix()
		request.CreatedAtFrom = &createdAtFromUnix
	}
	if !criteria.CreatedAtTo.IsZero() {
		createdAtToUnix := criteria.CreatedAtTo.UTC().Unix()
		request.CreatedAtTo = &createdAtToUnix
	}
	return request
}

func reportOperationFromProto(op *pbEngine.ReportOperation) model.ReportOperation {
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
