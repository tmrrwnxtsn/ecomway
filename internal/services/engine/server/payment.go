package server

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) CreatePayment(ctx context.Context, request *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	method, err := s.methodService.GetOne(ctx, model.TransactionTypePayment, request.GetCurrency(), request.GetExternalSystem(), request.GetExternalMethod())
	if err != nil {
		return nil, err
	}

	if err = s.limitService.ValidateAmount(request.GetAmount(), request.GetCurrency(), method); err != nil {
		return nil, err
	}

	return nil, nil
}
