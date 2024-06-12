package engine

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) AddToFavorites(ctx context.Context, data model.FavoritesData) error {
	request := &pb.FavoritesRequest{
		Type:           convert.OperationTypeToProto(data.OperationType),
		UserId:         data.UserID,
		Currency:       data.Currency,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
	}

	_, err := c.client.AddToFavorites(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return perr
		}
		return err
	}

	return nil
}

func (c *Client) RemoveFromFavorites(ctx context.Context, data model.FavoritesData) error {
	request := &pb.FavoritesRequest{
		Type:           convert.OperationTypeToProto(data.OperationType),
		UserId:         data.UserID,
		Currency:       data.Currency,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
	}

	_, err := c.client.RemoveFromFavorites(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return perr
		}
		return err
	}

	return nil
}
