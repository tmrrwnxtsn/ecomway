package integration

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) CreatePayout(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error) {
	var result model.CreatePayoutResult

	pbAdditional, err := structpb.NewStruct(data.AdditionalData)
	if err != nil {
		return result, err
	}

	request := &pb.CreatePayoutRequest{
		OperationId:    data.OperationID,
		UserId:         data.UserID,
		LangCode:       data.LangCode,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
		Amount:         data.Amount,
		Currency:       data.Currency,
		AdditionalData: pbAdditional,
	}

	if data.Tool != nil {
		request.Tool = convert.ToolToProto(data.Tool)
	}

	response, err := c.client.CreatePayout(ctx, request)
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

	return result, nil
}
