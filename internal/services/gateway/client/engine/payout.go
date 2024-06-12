package engine

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
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
		if perr := perror.FromProto(err); perr != nil {
			return result, perr
		}
		return result, err
	}

	result.OperationID = response.GetOperationId()
	result.Status = model.OperationStatus(response.GetStatus())

	return result, nil
}

func (c *Client) ConfirmPayout(ctx context.Context, data model.ConfirmPayoutData) error {
	request := &pb.ConfirmPayoutRequest{
		OperationId:      data.OperationID,
		UserId:           data.UserID,
		ConfirmationCode: data.ConfirmationCode,
		LangCode:         data.LangCode,
	}

	_, err := c.client.ConfirmPayout(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return perr
		}
		return err
	}

	return nil
}

func (c *Client) ResendCode(ctx context.Context, opID, userID int64, langCode string) error {
	request := &pb.ResendConfirmationCodeRequest{
		OperationId: opID,
		UserId:      userID,
		LangCode:    langCode,
	}

	_, err := c.client.ResendConfirmationCode(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return perr
		}
		return err
	}

	return nil
}
