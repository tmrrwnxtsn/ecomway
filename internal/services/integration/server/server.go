package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Integration interface {
	AvailableMethods(ctx context.Context, opType model.OperationType, currency string) ([]model.Method, error)
	CreatePayment(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error)
	GetOperationStatus(ctx context.Context, data model.GetOperationStatusData) (model.GetOperationStatusResult, error)
	CreatePayout(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error)
}

type Server struct {
	server       *grpc.Server
	listener     net.Listener
	integrations map[string]Integration
	pb.UnimplementedIntegrationServiceServer
}

type Options struct {
	Server       *grpc.Server
	Listener     net.Listener
	Integrations map[string]Integration
}

func NewServer(opts Options) *Server {
	var s Server
	s.server = opts.Server
	s.listener = opts.Listener
	s.integrations = opts.Integrations
	return &s
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Close() error {
	s.server.GracefulStop()
	return s.listener.Close()
}
