package v1

import (
	"context"
	"reflect"
	"strings"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/tmrrwnxtsn/ecomway/api/swagger/report/v1" // generated by Swag CLI, you have to import it
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/middleware"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type OperationService interface {
	ReportOperations(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error)
	GetExternalOperationStatus(ctx context.Context, id int64) (model.OperationExternalStatus, error)
}

type SortingService interface {
	SortReportOperations(items []model.ReportOperation, orderField, orderType string) []model.ReportOperation
}

type SummaryService interface {
	CalculateReportOperationsSummary(items []model.ReportOperation) (totalAmount float64, totalCount int64)
}

type ToolService interface {
	AllTools(ctx context.Context, userID int64) ([]*model.Tool, error)
	RecoverTool(ctx context.Context, id string, userID int64, externalMethod string) error
	RemoveTool(ctx context.Context, id string, userID int64, externalMethod string) error
}

type Translator interface {
	Translate(lang, key string, args ...any) string
}

type Handler struct {
	operationService OperationService
	sortingService   SortingService
	summaryService   SummaryService
	toolService      ToolService
	translator       Translator
	validate         *validator.Validate
	apiKey           string
}

type HandlerOptions struct {
	OperationService OperationService
	SortingService   SortingService
	SummaryService   SummaryService
	ToolService      ToolService
	Translator       Translator
	APIKey           string
}

// NewHandler godoc
//
//	@title						Шлюз финансовой отчетности E-commerce системы
//	@version					1.0
//
//	@contact.name				Курмыза Павел
//	@contact.email				tmrrwnxtsn@gmail.com
//
//	@host						localhost:8081
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Секретный ключ
func NewHandler(opts HandlerOptions) *Handler {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name != "-" && name != "" {
			return name
		}

		name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		if name != "-" && name != "" {
			return name
		}

		return ""
	})

	return &Handler{
		operationService: opts.OperationService,
		sortingService:   opts.SortingService,
		summaryService:   opts.SummaryService,
		toolService:      opts.ToolService,
		translator:       opts.Translator,
		validate:         validate,
		apiKey:           opts.APIKey,
	}
}

func (h *Handler) Init(router fiber.Router) {
	router.Get("/swagger/*", swagger.HandlerDefault)

	apiV1 := router.Group("/api/v1")

	{
		apiV1.Use(recover.New())
		apiV1.Use(h.authorizationMiddleware())
	}

	{
		operations := apiV1.Group("/operation")
		{
			operations.Get("", h.operationList)
			operations.Get("/:id/external-status", h.operationExternalStatus).Use(middleware.NewAccessLog())
			operations.Put("/:id/change-status", h.operationChangeStatus).Use(middleware.NewAccessLog())
		}
	}

	{
		tools := apiV1.Group("/tool")
		{
			tools.Get("", h.toolList)
			tools.Put("/recover", h.toolRecover).Use(middleware.NewAccessLog())
			tools.Delete("/delete", h.toolDelete).Use(middleware.NewAccessLog())
		}
	}
}
