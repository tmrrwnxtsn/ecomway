package server

import (
	"context"
	"sync"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/api/proto/shared"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Server) AvailableMethods(ctx context.Context, request *pb.AvailableMethodsRequest) (*pb.AvailableMethodsResponse, error) {
	userID := request.GetUserId()
	txType := convert.TransactionTypeFromProto(request.GetTransactionType())

	var wg sync.WaitGroup

	methodsChan := make(chan []model.Method)

	for _, integration := range s.integrations {
		if integration != nil {
			wg.Add(1)
			go func(ctx context.Context, integration Integration, userID int64, txType model.TransactionType) {
				defer wg.Done()

				methods, err := integration.AvailableMethods(ctx, userID, txType)
				if err != nil {
					return
				}

				methodsChan <- methods
			}(ctx, integration, userID, txType)
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
