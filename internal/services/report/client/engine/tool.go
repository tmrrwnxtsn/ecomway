package engine

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
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
