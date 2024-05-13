package server

import (
	"context"

	pbEngine "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
)

func (s *Server) AvailableTools(ctx context.Context, request *pbEngine.AvailableToolsRequest) (*pbEngine.AvailableToolsResponse, error) {
	tools, err := s.toolService.All(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	pbTools := make([]*pb.Tool, 0, len(tools))
	for _, tool := range tools {
		if tool != nil {
			pbTools = append(pbTools, convert.ToolToProto(tool))
		}
	}

	return &pbEngine.AvailableToolsResponse{
		Tools: pbTools,
	}, nil
}

func (s *Server) EditTool(ctx context.Context, request *pbEngine.EditToolRequest) (*pbEngine.EditToolResponse, error) {
	id := request.GetId()
	userID := request.GetUserId()
	externalMethod := request.GetExternalMethod()
	name := request.GetName()

	edited, err := s.toolService.EditOne(ctx, id, userID, externalMethod, name)
	if err != nil {
		return nil, err
	}

	return &pbEngine.EditToolResponse{
		Tool: convert.ToolToProto(edited),
	}, nil
}
