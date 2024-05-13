package engine

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (c *Client) AvailableTools(ctx context.Context, userID int64) ([]*model.Tool, error) {
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

func (c *Client) EditTool(ctx context.Context, id string, userID int64, externalMethod, name string) (*model.Tool, error) {
	request := &pb.EditToolRequest{
		Id:             id,
		UserId:         userID,
		ExternalMethod: externalMethod,
		Name:           name,
	}

	response, err := c.client.EditTool(ctx, request)
	if err != nil {
		if perr := perror.FromProto(err); perr != nil {
			return nil, perr
		}
		return nil, err
	}

	tool := convert.ToolFromProto(response.GetTool())

	return tool, nil
}
