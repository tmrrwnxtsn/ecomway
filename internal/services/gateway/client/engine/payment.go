package engine

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) CreatePayment(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error) {
	var result model.CreatePaymentResult

	pbAdditional, err := structpb.NewStruct(data.AdditionalData)
	if err != nil {
		return result, err
	}

	request := &pb.CreatePaymentRequest{
		UserId:         data.UserID,
		LangCode:       data.LangCode,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
		Amount:         data.Amount,
		Currency:       data.Currency,
		AdditionalData: pbAdditional,
		ReturnUrls:     convert.ReturnURLsToProto(data.ReturnURLs),
	}

	if data.ToolID != "" {
		request.ToolId = &data.ToolID
	}

	response, err := c.client.CreatePayment(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return result, perr
		}
		return result, err
	}

	result.RedirectURL = response.GetRedirectUrl()
	result.OperationID = response.GetOperationId()
	result.Status = model.OperationStatus(response.GetStatus())

	return result, nil
}
