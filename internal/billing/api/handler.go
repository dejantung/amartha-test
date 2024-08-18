package api

import (
	"billing-engine/internal/billing/model"
	"billing-engine/internal/billing/service"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type BillingHandler struct {
	BillingService service.BillingServiceProvider
	log            logger.Logger
}

func (s *BillingHandler) CreateLoanHandler(c echo.Context) error {
	ctx := c.Request().Context()

	payload := model.CreateLoanPayload{}
	if err := c.Bind(&payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return err
	}

	result, err := s.BillingService.CreateLoan(ctx, payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}

func (s *BillingHandler) GetPaymentScheduleHandler(c echo.Context) error {
	ctx := c.Request().Context()

	payload := model.GetSchedulePayload{}
	if err := c.Bind(&payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return err
	}

	result, err := s.BillingService.GetPaymentSchedule(ctx, payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}

func (s *BillingHandler) IsCustomerDelinquentHandler(c echo.Context) error {
	ctx := c.Request().Context()

	customerID := c.Param("customer_id")
	// convert string to uuid
	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(400, "invalid customer id"))
	}

	result, err := s.BillingService.IsCustomerDelinquency(ctx, customerUUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}

func (s *BillingHandler) GetOutstandingBalanceHandler(c echo.Context) error {
	ctx := c.Request().Context()

	customerID := c.Param("customer_id")
	// convert string to uuid
	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(400, "invalid customer id"))
	}

	result, err := s.BillingService.GetOutstandingBalance(ctx, customerUUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}

func (s *BillingHandler) CreateCustomerHandler(c echo.Context) error {
	ctx := c.Request().Context()

	payload := model.CreateCustomerPayload{}
	if err := c.Bind(&payload); err != nil {
		return err
	}

	if err := c.Validate(payload); err != nil {
		return err
	}

	result, err := s.BillingService.CreateCustomer(ctx, payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}

func (s *BillingHandler) GetCustomerHandler(c echo.Context) error {
	ctx := c.Request().Context()

	result, err := s.BillingService.GetCustomer(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}

func NewBillingHandler(svc service.BillingServiceProvider) *BillingHandler {
	return &BillingHandler{
		BillingService: svc,
	}
}
