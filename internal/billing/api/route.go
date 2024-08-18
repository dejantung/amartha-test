package api

import "github.com/labstack/echo/v4"

func (s *BillingHandler) AddRoutes(e *echo.Echo) {
	loanGroup := e.Group("/loan")
	loanGroup.POST("", s.CreateLoanHandler)
	loanGroup.GET("/schedule", s.GetPaymentScheduleHandler)

	customerGroup := e.Group("/customer")
	customerGroup.GET("/:customer_id/delinquent", s.IsCustomerDelinquentHandler)
	customerGroup.POST("", s.CreateCustomerHandler)
	customerGroup.GET("", s.GetCustomerHandler)
	customerGroup.GET("/:customer_id/outstanding", s.GetOutstandingBalanceHandler)
}
