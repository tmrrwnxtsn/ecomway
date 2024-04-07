package app

import (
	"context"
	"io"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbEngine "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	pbIntegration "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/client/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/server"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/limit"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/method"
)

type App struct {
	srv     *server.Server
	closers []io.Closer
}

func New(configPath string) *App {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	grpcListener, err := net.Listen("tcp", cfg.Engine.GRPCAddress)
	if err != nil {
		log.Fatalf("listening tcp on %v: %v", cfg.Engine.GRPCAddress, err)
	}

	integrationConn, err := grpc.Dial(cfg.Services.Integration.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connecting integration service: %v", err)
	}

	integrationClient := integration.NewClient(pbIntegration.NewIntegrationServiceClient(integrationConn))

	methodService := method.NewService(integrationClient)
	limitService := limit.NewService()

	grpcServer := grpc.NewServer()
	srv := server.NewServer(server.Options{
		Server:        grpcServer,
		Listener:      grpcListener,
		MethodService: methodService,
		LimitService:  limitService,
	})
	pbEngine.RegisterEngineServiceServer(grpcServer, srv)

	return &App{
		srv:     srv,
		closers: []io.Closer{srv, integrationConn},
	}
}

func (a *App) Start() {
	var group errgroup.Group

	group.Go(func() error {
		return a.srv.Serve()
	})

	if err := group.Wait(); err != nil {
		log.Fatalf("app: %v", err)
	}
}

func (a *App) Stop(_ context.Context) {
	for _, c := range a.closers {
		if err := c.Close(); err != nil {
			log.Printf("failed to close resource: %v", err)
		}
	}
}
