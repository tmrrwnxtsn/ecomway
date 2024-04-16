package integration

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) CreatePayment(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error) {
	var result model.CreatePaymentResult

	pbAdditional, err := structpb.NewStruct(data.AdditionalData)
	if err != nil {
		return result, err
	}

	request := &pb.CreatePaymentRequest{
		OperationId:    data.OperationID,
		UserId:         data.UserID,
		LangCode:       data.LangCode,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
		Amount:         data.Amount,
		Currency:       data.Currency,
		AdditionalData: pbAdditional,
		ReturnUrls:     convert.ReturnURLsToProto(data.ReturnURLs),
	}

	if data.Tool != nil {
		request.Tool = convert.ToolToProto(data.Tool)
	}

	response, err := c.client.CreatePayment(ctx, request)
	if err != nil {
		return result, err
	}

	result.ExternalStatus = convert.OperationExternalStatusFromProto(response.GetExternalStatus())

	if response.RedirectUrl != nil {
		result.RedirectURL = response.GetRedirectUrl()
	}

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
