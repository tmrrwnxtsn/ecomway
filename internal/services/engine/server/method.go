package server

import (
	"context"
	"fmt"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	pbShared "github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) AvailableMethods(ctx context.Context, request *pb.AvailableMethodsRequest) (*pb.AvailableMethodsResponse, error) {
	var txType model.TransactionType
	switch request.GetTransactionType() {
	case pbShared.TransactionType_PAYMENT:
		txType = model.TransactionTypePayment
	case pbShared.TransactionType_PAYOUT:
		txType = model.TransactionTypePayout
	default:
		return nil, fmt.Errorf("unresolved transaction type: %q", request.GetTransactionType().String())
	}

	methods, err := s.methodService.AvailableMethods(ctx, request.GetUserId(), txType)
	if err != nil {
		return nil, err
	}

	return &pb.AvailableMethodsResponse{
		Methods: convert.MethodsToProto(methods),
	}, nil
}
