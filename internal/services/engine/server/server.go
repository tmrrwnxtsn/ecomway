package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type MethodService interface {
	All(ctx context.Context, txType model.TransactionType, currency string) ([]model.Method, error)
	GetOne(ctx context.Context, txType model.TransactionType, currency, externalSystem, externalMethod string) (*model.Method, error)
}

type LimitService interface {
	ValidateAmount(amount int64, currency string, method *model.Method) error
}

type Server struct {
	server        *grpc.Server
	listener      net.Listener
	methodService MethodService
	limitService  LimitService
	pb.UnimplementedEngineServiceServer
}

type Options struct {
	Server        *grpc.Server
	Listener      net.Listener
	MethodService MethodService
	LimitService  LimitService
}

func NewServer(opts Options) *Server {
	var s Server
	s.server = opts.Server
	s.listener = opts.Listener
	s.methodService = opts.MethodService
	s.limitService = opts.LimitService
	return &s
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Close() error {
	s.server.GracefulStop()
	return s.listener.Close()
}
