package v1

import (
	"context"
	"reflect"
	"strings"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/tmrrwnxtsn/ecomway/api/swagger/gateway/v1" // generated by Swag CLI, you have to import it
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/middleware"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type MethodService interface {
	AvailableMethods(ctx context.Context, opType model.OperationType, userID string, currency string) ([]model.Method, error)
}

type PaymentService interface {
	Create(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error)
}

type ToolService interface {
	AvailableTools(ctx context.Context, userID string) ([]*model.Tool, error)
	AvailableToolsGroupedByMethod(ctx context.Context, userID string) (map[string][]*model.Tool, error)
	EditTool(ctx context.Context, id string, userID string, externalMethod, name string) (*model.Tool, error)
	RemoveTool(ctx context.Context, id string, userID string, externalMethod string) error
}

type PayoutService interface {
	Create(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error)
	Confirm(ctx context.Context, data model.ConfirmPayoutData) error
	ResendCode(ctx context.Context, opID int64, userID string, langCode string) error
}

type FavoritesService interface {
	Add(ctx context.Context, data model.FavoritesData) error
	Remove(ctx context.Context, data model.FavoritesData) error
}

type OperationService interface {
	ReportOperations(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error)
}

type SortingService interface {
	SortReportOperations(items []model.ReportOperation, orderField, orderType string) []model.ReportOperation
}

type SummaryService interface {
	CalculateReportOperationsSummary(items []model.ReportOperation) (totalAmount float64, totalCount int64)
}

type Translator interface {
	Translate(lang, key string, args ...any) string
}

type Handler struct {
	methodService    MethodService
	paymentService   PaymentService
	toolService      ToolService
	payoutService    PayoutService
	favoritesService FavoritesService
	operationService OperationService
	sortingService   SortingService
	summaryService   SummaryService
	translator       Translator
	validate         *validator.Validate
	apiKey           string
}

type HandlerOptions struct {
	MethodService    MethodService
	PaymentService   PaymentService
	ToolService      ToolService
	PayoutService    PayoutService
	FavoritesService FavoritesService
	OperationService OperationService
	SortingService   SortingService
	SummaryService   SummaryService
	Translator       Translator
	APIKey           string
}

// NewHandler godoc
//
//	@title						Платежный шлюз для E-commerce системы
//	@version					1.0
//
//	@contact.name				Курмыза Павел
//	@contact.email				tmrrwnxtsn@gmail.com
//
//	@host						localhost:8080
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
		methodService:    opts.MethodService,
		paymentService:   opts.PaymentService,
		toolService:      opts.ToolService,
		payoutService:    opts.PayoutService,
		favoritesService: opts.FavoritesService,
		operationService: opts.OperationService,
		sortingService:   opts.SortingService,
		summaryService:   opts.SummaryService,
		translator:       opts.Translator,
		validate:         validate,
		apiKey:           opts.APIKey,
	}
}

func (h *Handler) Init(router fiber.Router) {
	router.Get("/swagger/*", swagger.HandlerDefault)

	apiV1 := router.Group("/api/v1")

	{
		apiV1.Use(middleware.NewAccessLog())
		apiV1.Use(recover.New())
		apiV1.Use(h.authorizationMiddleware())
	}

	{
		payment := apiV1.Group("/payment")
		{
			payment.Get("/methods", h.paymentMethods)
			payment.Post("/create", h.paymentCreate)
		}

		payout := apiV1.Group("/payout")
		{
			payout.Get("/methods", h.payoutMethods)
			payout.Post("/create", h.payoutCreate)
			payout.Put("/:id/confirm", h.payoutConfirm)
			payout.Put("/:id/resend-code", h.payoutResendCode)
		}

		tools := apiV1.Group("/tool")
		{
			tools.Get("", h.toolList)
			tools.Put("/edit", h.toolEdit)
			tools.Delete("/remove", h.toolRemove)
		}

		favorites := apiV1.Group("/favorites/:operation_type")
		{
			favorites.Post("", h.favoritesAdd)
			favorites.Delete("", h.favoritesRemove)
		}

		operations := apiV1.Group("/operation")
		{
			operations.Get("", h.operationList)
		}
	}
}
