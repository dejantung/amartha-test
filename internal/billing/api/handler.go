package api

import (
	"billing-engine/internal/billing/model"
	"billing-engine/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *AppServer) CreateLoanHandler(c echo.Context) error {
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

func (s *AppServer) GetPaymentScheduleHandler(c echo.Context) error {
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
