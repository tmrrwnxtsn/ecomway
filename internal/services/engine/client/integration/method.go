package integration

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) AllAvailableMethods(ctx context.Context, opType model.OperationType, currency string) ([]model.Method, error) {
	request := &pb.AvailableMethodsRequest{
		OperationType: convert.OperationTypeToProto(opType),
		Currency:      currency,
	}

	return c.availableMethods(ctx, request)
}

func (c *Client) AvailableMethodsByExternalSystem(ctx context.Context, opType model.OperationType, currency, externalSystem string) ([]model.Method, error) {
	request := &pb.AvailableMethodsRequest{
		OperationType:  convert.OperationTypeToProto(opType),
		Currency:       currency,
		ExternalSystem: &externalSystem,
	}

	return c.availableMethods(ctx, request)
}

func (c *Client) availableMethods(ctx context.Context, request *pb.AvailableMethodsRequest) ([]model.Method, error) {
	response, err := c.client.AvailableMethods(ctx, request)
	if err != nil {
		return nil, err
	}

	result := convert.MethodsFromProto(response.GetMethods())

	return result, nil
}
