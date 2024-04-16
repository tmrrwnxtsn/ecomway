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
