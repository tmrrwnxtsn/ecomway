package server

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
)

func (s *Server) AvailableMethods(ctx context.Context, request *pb.AvailableMethodsRequest) (*pb.AvailableMethodsResponse, error) {
	opType := convert.OperationTypeFromProto(request.GetOperationType())

	methods, err := s.methodService.All(ctx, opType, request.GetCurrency())
	if err != nil {
		return nil, err
	}

	if err = s.favoritesService.FillForMethods(ctx, opType, request.GetUserId(), methods); err != nil {
		return nil, err
	}

	return &pb.AvailableMethodsResponse{
		Methods: convert.MethodsToProto(methods),
	}, nil
}
