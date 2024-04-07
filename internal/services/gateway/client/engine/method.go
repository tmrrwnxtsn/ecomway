package engine

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) AvailableMethods(ctx context.Context, txType model.TransactionType, userID int64, currency string) ([]model.Method, error) {
	request := &pb.AvailableMethodsRequest{
		TransactionType: convert.TransactionTypeToProto(txType),
		Currency:        currency,
		UserId:          userID,
	}

	response, err := c.client.AvailableMethods(ctx, request)
	if err != nil {
		return nil, err
	}

	result := convert.MethodsFromProto(response.GetMethods())

	return result, nil
}
