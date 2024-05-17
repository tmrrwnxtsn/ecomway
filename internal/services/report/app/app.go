package app

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbEngine "github.com/tmrrwnxtsn/ecomway/api/proto/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/translate"
	"github.com/tmrrwnxtsn/ecomway/internal/services/report/api"
	"github.com/tmrrwnxtsn/ecomway/internal/services/report/api/v1"
	"github.com/tmrrwnxtsn/ecomway/internal/services/report/client/engine"
	"github.com/tmrrwnxtsn/ecomway/internal/services/report/config"
	"github.com/tmrrwnxtsn/ecomway/internal/services/report/service/operation"
)

type App struct {
	app     *fiber.App
	config  config.ReportConfig
	closers []io.Closer
}

func New(configPath string) *App {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	engineConn, err := grpc.Dial(cfg.Services.Engine.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connecting engine service: %v", err)
	}

	engineClient := engine.NewClient(pbEngine.NewEngineServiceClient(engineConn))

	operationService := operation.NewService(engineClient)
	translator := translate.NewTranslator("en", "ru")

	apiHandlerV1 := v1.NewHandler(v1.HandlerOptions{
		OperationService: operationService,
		Translator:       translator,
		APIKey:           cfg.Report.APIKey,
	})
	apiServer := api.NewServer(apiHandlerV1)

	app := fiber.New()
	apiServer.Init(app)

	return &App{
		app:     app,
		config:  cfg.Report,
		closers: []io.Closer{engineConn},
	}
}

func (a *App) Start() {
	var group errgroup.Group

	group.Go(func() error {
		return a.app.Listen(a.config.HTTPAddress)
	})

	if err := group.Wait(); err != nil {
		log.Fatalf("app: %v", err)
	}
}

func (a *App) Stop(_ context.Context) {
	if err := a.app.Shutdown(); err != nil {
		log.Printf("failed to shutdown the app: %v", err)
	}

	for _, c := range a.closers {
		if err := c.Close(); err != nil {
			log.Printf("failed to close resource: %v", err)
		}
	}
}
