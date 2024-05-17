package engine

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) Operations(ctx context.Context, userID int64) ([]model.Operation, error) {
	request := &pb.GetOperationsRequest{
		UserId: userID,
	}

	response, err := c.client.GetOperations(ctx, request)
	if err != nil {
		return nil, err
	}

	operations := make([]model.Operation, 0, len(response.GetOperations()))
	for _, pbOp := range response.GetOperations() {
		operations = append(operations, convert.OperationFromProto(pbOp))
	}
	return operations, nil
}
