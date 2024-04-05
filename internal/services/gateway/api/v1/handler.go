package v1

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type MethodService interface {
	AvailableMethods(ctx context.Context, userID int64, txType model.TransactionType) ([]model.Method, error)
}

type Handler struct {
	methodService MethodService
	validate      *validator.Validate
	apiKey        string
}

type HandlerOptions struct {
	MethodService MethodService
	APIKey        string
}

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
		methodService: opts.MethodService,
		validate:      validate,
		apiKey:        opts.APIKey,
	}
}

func (h *Handler) Init(router fiber.Router) {
	router.Use(h.authorizationMiddleware())

	apiV1 := router.Group("/api/v1")
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
			payout.Put("/confirm", h.payoutConfirm)
			payout.Put("/resendCode", h.payoutResendCode)
		}
	}
}
