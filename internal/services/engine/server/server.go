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
}

type ToolService interface {
	All(ctx context.Context, userID int64) ([]*model.Tool, error)
}

type PayoutService interface {
	Create(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error)
}

type Server struct {
	server         *grpc.Server
	listener       net.Listener
	methodService  MethodService
	limitService   LimitService
	paymentService PaymentService
	toolService    ToolService
	payoutService  PayoutService
	pb.UnimplementedEngineServiceServer
}

type Options struct {
	Server         *grpc.Server
	Listener       net.Listener
	MethodService  MethodService
	LimitService   LimitService
	PaymentService PaymentService
	ToolService    ToolService
	PayoutService  PayoutService
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
	return &s
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Close() error {
	s.server.GracefulStop()
	return s.listener.Close()
}
