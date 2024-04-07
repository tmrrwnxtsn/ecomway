package server

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
)

func (s *Server) AvailableMethods(ctx context.Context, request *pb.AvailableMethodsRequest) (*pb.AvailableMethodsResponse, error) {
	txType := convert.TransactionTypeFromProto(request.GetTransactionType())

	methods, err := s.methodService.All(ctx, txType, request.GetCurrency())
	if err != nil {
		return nil, err
	}

	return &pb.AvailableMethodsResponse{
		Methods: convert.MethodsToProto(methods),
	}, nil
}
