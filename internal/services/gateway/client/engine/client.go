package engine

import (
	"context"
	"fmt"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	pbShared "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Client struct {
	client pb.EngineServiceClient
}

func NewClient(client pb.EngineServiceClient) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) AvailableMethods(ctx context.Context, userID int64, txType model.TransactionType) ([]model.Method, error) {
	var pbTxType pbShared.TransactionType
	switch txType {
	case model.TransactionTypePayment:
		pbTxType = pbShared.TransactionType_PAYMENT
	case model.TransactionTypePayout:
		pbTxType = pbShared.TransactionType_PAYOUT
	default:
		return nil, fmt.Errorf("unresolved transaction type: %q", txType)
	}

	request := &pb.AvailableMethodsRequest{
		UserId:          userID,
		TransactionType: pbTxType,
	}

	response, err := c.client.AvailableMethods(ctx, request)
	if err != nil {
		return nil, err
	}

	result := convert.MethodsFromProto(response.GetMethods())

	return result, nil
}
