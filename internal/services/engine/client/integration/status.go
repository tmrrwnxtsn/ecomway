package integration

import (
	"context"
	"time"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) GetOperationStatus(ctx context.Context, data model.GetOperationStatusData) (model.GetOperationStatusResult, error) {
	var result model.GetOperationStatusResult

	request := &pb.GetOperationStatusRequest{
		OperationId:    data.OperationID,
		OperationType:  convert.OperationTypeToProto(data.OperationType),
		UserId:         data.UserID,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
		Amount:         data.Amount,
		Currency:       data.Currency,
		CreatedAt:      data.CreatedAt.UTC().Unix(),
	}

	if data.ExternalID != "" {
		request.ExternalId = &data.ExternalID
	}

	response, err := c.client.GetOperationStatus(ctx, request)
	if err != nil {
		return result, err
	}

	result.ExternalStatus = convert.OperationExternalStatusFromProto(response.GetExternalStatus())

	if response.ExternalId != nil {
		result.ExternalID = response.GetExternalId()
	}

	if response.ProcessedAt != nil {
		result.ProcessedAt = time.Unix(response.GetProcessedAt(), 0).UTC()
	}

	if response.FailReason != nil {
		result.FailReason = response.GetFailReason()
	}

	if response.NewAmount != nil {
		result.NewAmount = response.GetNewAmount()
	}

	if response.Tool != nil {
		result.Tool = convert.ToolFromProto(response.GetTool())
	}

	return result, nil
}
