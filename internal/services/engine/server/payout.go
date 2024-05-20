package server

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) CreatePayout(ctx context.Context, request *pb.CreatePayoutRequest) (*pb.CreatePayoutResponse, error) {
	method, err := s.methodService.GetOne(ctx, model.OperationTypePayout, request.GetCurrency(), request.GetExternalSystem(), request.GetExternalMethod())
	if err != nil {
		return nil, err
	}

	if err = s.limitService.ValidateAmount(request.GetAmount(), request.GetCurrency(), method); err != nil {
		return nil, err
	}

	data := model.CreatePayoutData{
		AdditionalData: request.GetAdditionalData().AsMap(),
		ExternalSystem: request.GetExternalSystem(),
		ExternalMethod: request.GetExternalMethod(),
		Currency:       request.GetCurrency(),
		LangCode:       request.GetLangCode(),
		UserID:         request.GetUserId(),
		ToolID:         request.GetToolId(),
		Amount:         request.GetAmount(),
	}

	result, err := s.payoutService.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePayoutResponse{
		OperationId: result.OperationID,
		Status:      string(result.Status),
	}, nil
}

func (s *Server) ConfirmPayout(ctx context.Context, request *pb.ConfirmPayoutRequest) (*emptypb.Empty, error) {
	data := model.ConfirmPayoutData{
		ConfirmationCode: request.GetConfirmationCode(),
		LangCode:         request.GetLangCode(),
		UserID:           request.GetUserId(),
		OperationID:      request.GetOperationId(),
	}

	if err := s.payoutService.Confirm(ctx, data); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
