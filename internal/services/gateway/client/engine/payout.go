package engine

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) CreatePayout(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error) {
	var result model.CreatePayoutResult

	pbAdditional, err := structpb.NewStruct(data.AdditionalData)
	if err != nil {
		return result, err
	}

	request := &pb.CreatePayoutRequest{
		UserId:         data.UserID,
		ToolId:         data.ToolID,
		LangCode:       data.LangCode,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
		Amount:         data.Amount,
		Currency:       data.Currency,
		AdditionalData: pbAdditional,
	}

	response, err := c.client.CreatePayout(ctx, request)
	if err != nil {
		return result, err
	}

	result.OperationID = response.GetOperationId()
	result.Status = model.OperationStatus(response.GetStatus())

	return result, nil
}
