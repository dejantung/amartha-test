package api

import (
	"billing-engine/internal/payment/model"
	"billing-engine/internal/payment/service"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PaymentHandler struct {
	BillingService service.PaymentServiceProvider
	log            logger.Logger
}

func NewPaymentHandler(s service.PaymentServiceProvider, l logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		BillingService: s,
		log:            l,
	}
}

func (s *PaymentHandler) ProcessPaymentHandler(c echo.Context) error {
	ctx := c.Request().Context()

	payload := model.ProcessPaymentPayload{}
	if err := c.Bind(&payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return err
	}

	result, err := s.BillingService.ProcessPayment(ctx, payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}
