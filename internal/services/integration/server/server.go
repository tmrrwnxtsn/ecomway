package server

import (
	"net"

	"google.golang.org/grpc"

	pb "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
)

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
