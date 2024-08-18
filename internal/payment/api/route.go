package api

import "github.com/labstack/echo/v4"

func (s *PaymentHandler) AddRoutes(e *echo.Echo) {
	paymentGroup := e.Group("/payment")
	paymentGroup.POST("", s.ProcessPaymentHandler)
}
