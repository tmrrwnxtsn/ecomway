package engine

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) AvailableTools(ctx context.Context, userID string) ([]*model.Tool, error) {
	request := &pb.AvailableToolsRequest{
		UserId: userID,
	}

	response, err := c.client.AvailableTools(ctx, request)
	if err != nil {
		return nil, err
	}

	tools := make([]*model.Tool, 0, len(response.GetTools()))
	for _, pbTool := range response.GetTools() {
		tools = append(tools, convert.ToolFromProto(pbTool))
	}
	return tools, nil
}

func (c *Client) RecoverTool(ctx context.Context, id string, userID string, externalMethod string) error {
	request := &pb.RecoverToolRequest{
		Id:             id,
		UserId:         userID,
		ExternalMethod: externalMethod,
	}

	_, err := c.client.RecoverTool(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return perr
		}
		return err
	}

	return nil
}

func (c *Client) RemoveTool(ctx context.Context, id string, userID string, externalMethod string) error {
	request := &pb.RemoveToolRequest{
		Id:             id,
		UserId:         userID,
		ExternalMethod: externalMethod,
		ActionSource:   pb.ActionSource_ACTION_SOURCE_ADMINISTRATOR,
	}

	_, err := c.client.RemoveTool(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return perr
		}
		return err
	}

	return nil
}
