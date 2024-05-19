package server

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pbEngine "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
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

func (s *Server) RemoveTool(ctx context.Context, request *pbEngine.RemoveToolRequest) (*emptypb.Empty, error) {
	id := request.GetId()
	userID := request.GetUserId()
	externalMethod := request.GetExternalMethod()
	actionSource := actionSourceFromProto(request.GetActionSource())

	if err := s.toolService.RemoveOne(ctx, id, userID, externalMethod, actionSource); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) RecoverTool(ctx context.Context, request *pbEngine.RecoverToolRequest) (*emptypb.Empty, error) {
	id := request.GetId()
	userID := request.GetUserId()
	externalMethod := request.GetExternalMethod()

	if err := s.toolService.RecoverOne(ctx, id, userID, externalMethod); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func actionSourceFromProto(actionSource pbEngine.ActionSource) model.ActionSource {
	switch actionSource {
	case pbEngine.ActionSource_ACTION_SOURCE_DEFAULT:
		return model.ActionSourceDefault
	case pbEngine.ActionSource_ACTION_SOURCE_ADMINISTRATOR:
		return model.ActionSourceAdministrator
	default:
		return model.ActionSourceDefault
	}
}
