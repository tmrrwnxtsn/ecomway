package server

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pbEngine "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) AddToFavorites(ctx context.Context, request *pbEngine.FavoritesRequest) (*emptypb.Empty, error) {
	opType := convert.OperationTypeFromProto(request.GetType())
	currency := request.GetCurrency()
	externalSystem := request.GetExternalSystem()
	externalMethod := request.GetExternalMethod()

	_, err := s.methodService.GetOne(ctx, opType, currency, externalSystem, externalMethod)
	if err != nil {
		return nil, err
	}

	data := model.FavoritesData{
		OperationType:  opType,
		Currency:       currency,
		ExternalSystem: externalSystem,
		ExternalMethod: externalMethod,
		UserID:         request.GetUserId(),
	}

	if err = s.favoritesService.AddToFavorites(ctx, data); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) RemoveFromFavorites(ctx context.Context, request *pbEngine.FavoritesRequest) (*emptypb.Empty, error) {
	opType := convert.OperationTypeFromProto(request.GetType())
	currency := request.GetCurrency()
	externalSystem := request.GetExternalSystem()
	externalMethod := request.GetExternalMethod()

	_, err := s.methodService.GetOne(ctx, opType, currency, externalSystem, externalMethod)
	if err != nil {
		return nil, err
	}

	data := model.FavoritesData{
		OperationType:  opType,
		Currency:       currency,
		ExternalSystem: externalSystem,
		ExternalMethod: externalMethod,
		UserID:         request.GetUserId(),
	}

	if err = s.favoritesService.RemoveFromFavorites(ctx, data); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
