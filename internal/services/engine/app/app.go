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
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/client/smtp"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/migrator"
	oprepo "github.com/tmrrwnxtsn/ecomway/internal/services/engine/repository/operation"
	toolrepo "github.com/tmrrwnxtsn/ecomway/internal/services/engine/repository/tool"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/repository/user"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/scheduler"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/server"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/favorites"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/limit"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/method"
	opservice "github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/operation"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/payment"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/payout"
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

	env := cfg.Engine.Environment

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

	integrationClient := integration.NewClient(pbIntegration.NewIntegrationServiceClient(integrationConn))
	smtpClient := smtp.NewClient(cfg.Services.SMTP.Host, cfg.Services.SMTP.Port, cfg.Services.SMTP.Username, cfg.Services.SMTP.Password)
	if err = smtpClient.Connect(); err != nil {
		log.Fatalf("connecting to SMTP server: %v", err)
	}

	operationRepository := oprepo.NewRepository(postgresConn)
	toolRepository := toolrepo.NewRepository(postgresConn)
	userRepository := user.NewRepository(postgresConn)

	methodService := method.NewService(integrationClient)
	limitService := limit.NewService()
	paymentService := payment.NewService(operationRepository, integrationClient, toolRepository)
	toolService := toolservice.NewService(toolRepository)
	payoutService := payout.NewService(
		operationRepository,
		integrationClient,
		toolRepository,
		smtpClient,
		cfg.Engine.WrongConfirmationCodeLimit,
		env == "dev",
	)
	operationService := opservice.NewService(operationRepository)
	favoritesService := favorites.NewService(userRepository)

	if cfg.Engine.Scheduler.IsEnabled {
		var tasks []scheduler.BackgroundTask

		if cfg.Engine.Scheduler.Tasks.FinalizeOperations.IsEnabled {
			tasks = append(tasks, scheduler.NewFinalizeOperationsTask(
				cfg.Engine.Scheduler.Tasks.FinalizeOperations,
				operationService,
				integrationClient,
				paymentService,
				payoutService,
			))
		}
		if cfg.Engine.Scheduler.Tasks.RequestPayouts.IsEnabled {
			tasks = append(tasks, scheduler.NewRequestPayoutsTask(
				cfg.Engine.Scheduler.Tasks.RequestPayouts,
				operationService,
				payoutService,
			))
		}

		scheduler.NewScheduler(tasks...).Start(ctx)
	}

	grpcServer := grpc.NewServer()
	srv := server.NewServer(server.Options{
		Server:            grpcServer,
		Listener:          grpcListener,
		MethodService:     methodService,
		LimitService:      limitService,
		PaymentService:    paymentService,
		ToolService:       toolService,
		PayoutService:     payoutService,
		OperationService:  operationService,
		FavoritesService:  favoritesService,
		IntegrationClient: integrationClient,
	})
	pbEngine.RegisterEngineServiceServer(grpcServer, srv)

	return &App{
		srv:     srv,
		storage: postgresConn,
		closers: []io.Closer{srv, integrationConn, smtpClient},
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
