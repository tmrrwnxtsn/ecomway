package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type MethodService interface {
	AvailableMethods(ctx context.Context, userID int64, txType model.TransactionType) ([]model.Method, error)
}

type Server struct {
	server        *grpc.Server
	listener      net.Listener
	methodService MethodService
	pb.UnimplementedEngineServiceServer
}

type Options struct {
	Server        *grpc.Server
	Listener      net.Listener
	MethodService MethodService
}

func NewServer(opts Options) *Server {
	var s Server
	s.server = opts.Server
	s.listener = opts.Listener
	s.methodService = opts.MethodService
	return &s
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Close() error {
	s.server.GracefulStop()
	return s.listener.Close()
}
