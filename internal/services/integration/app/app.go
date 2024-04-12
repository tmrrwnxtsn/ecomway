package app

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net"
	"os"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pbIntegration "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/server"
)

type App struct {
	srv     *server.Server
	closers []io.Closer
}

func New(configPath string) *App {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	grpcListener, err := net.Listen("tcp", cfg.Integration.GRPCAddress)
	if err != nil {
		log.Fatalf("listening tcp on %v: %v", cfg.Integration.GRPCAddress, err)
	}

	integrations := map[string]server.Integration{
		yookassa.ExternalSystem: yookassa.NewIntegration(cfg.Integration.YooKassa),
	}

	grpcServer := grpc.NewServer()
	srv := server.NewServer(server.Options{
		Server:       grpcServer,
		Listener:     grpcListener,
		Integrations: integrations,
	})
	pbIntegration.RegisterIntegrationServiceServer(grpcServer, srv)

	return &App{
		srv:     srv,
		closers: []io.Closer{srv},
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
