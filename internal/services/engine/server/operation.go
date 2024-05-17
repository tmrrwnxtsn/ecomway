package server

import (
	"context"

	pbEngine "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	pb "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) GetOperations(ctx context.Context, request *pbEngine.GetOperationsRequest) (*pbEngine.GetOperationsResponse, error) {
	criteria := model.OperationCriteria{}

	operations, err := s.operationService.All(ctx, criteria)
	if err != nil {
		return nil, err
	}

	pbOperations := make([]*pb.Operation, 0, len(operations))
	for _, op := range operations {
		if op != nil {
			pbOperations = append(pbOperations, convert.OperationToProto(op))
		}
	}

	return &pbEngine.GetOperationsResponse{
		Operations: pbOperations,
	}, nil
}
