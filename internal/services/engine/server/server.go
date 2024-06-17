package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type MethodService interface {
	All(ctx context.Context, opType model.OperationType, currency string) ([]model.Method, error)
	GetOne(ctx context.Context, opType model.OperationType, currency, externalSystem, externalMethod string) (*model.Method, error)
}

type LimitService interface {
	ValidateAmount(amount int64, currency string, method *model.Method) error
}

type PaymentService interface {
	Create(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error)
	Success(ctx context.Context, data model.SuccessPaymentData) error
	Fail(ctx context.Context, data model.FailPaymentData) error
}

type ToolService interface {
	All(ctx context.Context, userID string) ([]*model.Tool, error)
	EditOne(ctx context.Context, id, userID, externalMethod, name string) (*model.Tool, error)
	RemoveOne(ctx context.Context, id, userID, externalMethod string, source model.ActionSource) error
	RecoverOne(ctx context.Context, id, userID, externalMethod string) error
	GetOne(ctx context.Context, id, userID, externalMethod string) (*model.Tool, error)
}

type PayoutService interface {
	Create(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error)
	Confirm(ctx context.Context, data model.ConfirmPayoutData) error
	ResendCode(ctx context.Context, opID int64, userID, langCode string) error
	Fail(ctx context.Context, data model.FailPayoutData) error
	Success(ctx context.Context, data model.SuccessPayoutData) error
}

type OperationService interface {
	AllForReport(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error)
	GetOne(ctx context.Context, criteria model.OperationCriteria) (*model.Operation, error)
	ChangeStatus(ctx context.Context, id int64, newStatus model.OperationStatus, newExternalStatus model.OperationExternalStatus) (model.OperationChangeStatusResult, error)
}

type FavoritesService interface {
	AddToFavorites(ctx context.Context, data model.FavoritesData) error
	RemoveFromFavorites(ctx context.Context, data model.FavoritesData) error
	FillForMethods(ctx context.Context, opType model.OperationType, userID string, methods []model.Method) error
}

type IntegrationClient interface {
	GetOperationStatus(ctx context.Context, data model.GetOperationStatusData) (model.GetOperationStatusResult, error)
}

type Server struct {
	server            *grpc.Server
	listener          net.Listener
	methodService     MethodService
	limitService      LimitService
	paymentService    PaymentService
	toolService       ToolService
	payoutService     PayoutService
	operationService  OperationService
	favoritesService  FavoritesService
	integrationClient IntegrationClient
	pb.UnimplementedEngineServiceServer
}

type Options struct {
	Server            *grpc.Server
	Listener          net.Listener
	MethodService     MethodService
	LimitService      LimitService
	PaymentService    PaymentService
	ToolService       ToolService
	PayoutService     PayoutService
	OperationService  OperationService
	FavoritesService  FavoritesService
	IntegrationClient IntegrationClient
}

func NewServer(opts Options) *Server {
	var s Server
	s.server = opts.Server
	s.listener = opts.Listener
	s.methodService = opts.MethodService
	s.limitService = opts.LimitService
	s.paymentService = opts.PaymentService
	s.toolService = opts.ToolService
	s.payoutService = opts.PayoutService
	s.operationService = opts.OperationService
	s.favoritesService = opts.FavoritesService
	s.integrationClient = opts.IntegrationClient
	return &s
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Close() error {
	s.server.GracefulStop()
	return s.listener.Close()
}
