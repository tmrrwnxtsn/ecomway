package engine

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) AvailableMethods(ctx context.Context, opType model.OperationType, userID int64, currency string) ([]model.Method, error) {
	request := &pb.AvailableMethodsRequest{
		OperationType: convert.OperationTypeToProto(opType),
		Currency:      currency,
		UserId:        userID,
	}

	response, err := c.client.AvailableMethods(ctx, request)
	if err != nil {
		return nil, err
	}

	result := convert.MethodsFromProto(response.GetMethods())

	return result, nil
}
