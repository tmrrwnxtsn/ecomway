package server

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) AvailableMethods(ctx context.Context, request *pb.AvailableMethodsRequest) (*pb.AvailableMethodsResponse, error) {
	opType := convert.OperationTypeFromProto(request.GetOperationType())

	if request.ExternalSystem != nil {
		integration, ok := s.integrations[request.GetExternalSystem()]
		if !ok || integration == nil {
			return nil, fmt.Errorf("unknown external system: %q", request.GetExternalSystem())
		}

		methods, err := integration.AvailableMethods(ctx, opType, request.GetCurrency())
		if err != nil {
			return nil, err
		}

		return &pb.AvailableMethodsResponse{
			Methods: convert.MethodsToProto(methods),
		}, nil
	}

	var wg sync.WaitGroup

	methodsChan := make(chan []model.Method)

	for _, integration := range s.integrations {
		if integration != nil {
			wg.Add(1)
			go func(ctx context.Context, integration Integration, opType model.OperationType, currency string) {
				defer wg.Done()

				methods, err := integration.AvailableMethods(ctx, opType, currency)
				if err != nil {
					return
				}

				methodsChan <- methods
			}(ctx, integration, opType, request.GetCurrency())
		}
	}

	go func() {
		wg.Wait()
		close(methodsChan)
	}()

	var pbMethods []*shared.Method
	for methods := range methodsChan {
		pbMethods = append(pbMethods, convert.MethodsToProto(methods)...)
	}

	return &pb.AvailableMethodsResponse{
		Methods: pbMethods,
	}, nil
}
