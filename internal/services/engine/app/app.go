package app

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbEngine "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	pbIntegration "github.com/tmrrwnxtsn/ecomway/api/proto/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/client/integration"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/migrator"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/repository/operation"
	toolrepo "github.com/tmrrwnxtsn/ecomway/internal/services/engine/repository/tool"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/scheduler"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/server"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/limit"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/method"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/payment"
	toolservice "github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/tool"
)

type App struct {
	srv     *server.Server
	storage *pgxpool.Pool
	closers []io.Closer
}

func New(configPath string) *App {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	postgresMigrator, err := migrator.NewPostgresMigrator(cfg.Engine.Storage.DatabaseURL)
	if err != nil {
		log.Fatalf("initializing migrator: %v", err)
	}
	if err = postgresMigrator.Migrate(); err != nil {
		log.Fatalf("applying migrations: %v", err)
	}
	if err = postgresMigrator.Close(); err != nil {
		log.Printf("closing migrator: %v", err)
	}

	postgresConn, err := pgxpool.Connect(ctx, cfg.Engine.Storage.DatabaseURL)
	if err != nil {
		log.Fatalf("connecting storage: %v", err)
	}

	grpcListener, err := net.Listen("tcp", cfg.Engine.GRPCAddress)
	if err != nil {
		log.Fatalf("listening tcp on %v: %v", cfg.Engine.GRPCAddress, err)
	}

	integrationConn, err := grpc.Dial(cfg.Services.Integration.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connecting integration service: %v", err)
	}

	operationRepository := operation.NewRepository(postgresConn)
	toolRepository := toolrepo.NewRepository(postgresConn)

	integrationClient := integration.NewClient(pbIntegration.NewIntegrationServiceClient(integrationConn))

	methodService := method.NewService(integrationClient)
	limitService := limit.NewService()
	paymentService := payment.NewService(operationRepository, integrationClient, toolRepository)
	toolService := toolservice.NewService(toolRepository)

	if cfg.Engine.Scheduler.IsEnabled {
		var tasks []scheduler.BackgroundTask

		if cfg.Engine.Scheduler.Tasks.FinalizeOperations.IsEnabled {
			tasks = append(tasks, scheduler.NewFinalizeOperationsTask(
				cfg.Engine.Scheduler.Tasks.FinalizeOperations,
				operationRepository,
				integrationClient,
				paymentService,
			))
		}

		scheduler.NewScheduler(tasks...).Start(ctx)
	}

	grpcServer := grpc.NewServer()
	srv := server.NewServer(server.Options{
		Server:         grpcServer,
		Listener:       grpcListener,
		MethodService:  methodService,
		LimitService:   limitService,
		PaymentService: paymentService,
		ToolService:    toolService,
	})
	pbEngine.RegisterEngineServiceServer(grpcServer, srv)

	return &App{
		srv:     srv,
		storage: postgresConn,
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

	a.storage.Close()
}
