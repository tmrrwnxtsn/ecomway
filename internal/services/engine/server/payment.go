package server

import (
	"context"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) CreatePayment(ctx context.Context, request *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	method, err := s.methodService.GetOne(ctx, model.OperationTypePayment, request.GetCurrency(), request.GetExternalSystem(), request.GetExternalMethod())
	if err != nil {
		return nil, err
	}

	if err = s.limitService.ValidateAmount(request.GetAmount(), request.GetCurrency(), method); err != nil {
		return nil, err
	}

	data := model.CreatePaymentData{
		ReturnURLs:     convert.ReturnURLsFromProto(request.GetReturnUrls()),
		AdditionalData: request.GetAdditionalData().AsMap(),
		ExternalSystem: request.GetExternalSystem(),
		ExternalMethod: request.GetExternalMethod(),
		Currency:       request.GetCurrency(),
		LangCode:       request.GetLangCode(),
		UserID:         request.GetUserId(),
		ToolID:         request.GetToolId(),
		Amount:         request.GetAmount(),
	}

	result, err := s.paymentService.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		RedirectUrl: result.RedirectURL,
		OperationId: result.OperationID,
	}, nil
}
