package model

import (
	"billing-engine/pkg/enum"
	"github.com/google/uuid"
)

type CreateLoanPayload struct {
	CustomerID uuid.UUID `json:"customer_id"`
	LoanAmount float64   `json:"loan_amount"`
}

type ScheduleResponse struct {
	ScheduleID     uuid.UUID          `json:"schedule_id"`
	LoanID         uuid.UUID          `json:"loan_id"`
	PaymentNo      int                `json:"payment_no"`
	PaymentDueDate string             `json:"payment_due_date"`
	PaymentAmount  float64            `json:"payment_amount"`
	PaymentStatus  enum.PaymentStatus `json:"payment_status"`
	IsMissPayment  bool               `json:"is_miss_payment"`
}

type CreateLoanResponse struct {
	LoanID     uuid.UUID          `json:"loan_id"`
	CustomerID uuid.UUID          `json:"customer_id"`
	LoanAmount float64            `json:"loan_amount"`
	Schedules  []ScheduleResponse `json:"schedules"`
}

type GetScheduleResponse struct {
	Schedules []ScheduleResponse `json:"schedules"`
}

type GetSchedulePayload struct {
	LoanID     uuid.UUID `query:"loan_id"`
	CustomerID uuid.UUID `query:"customer_id"`
}

type IsDelinquentPayload struct {
	LoanID     uuid.UUID `query:"loan_id"`
	CustomerID uuid.UUID `query:"customer_id"`
}

type IsDelinquentResponse struct {
	IsDelinquent bool `json:"is_delinquent"`
}

type GetCustomerResponse struct {
	CustomerID  uuid.UUID `json:"customer_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
}

type CreateCustomerPayload struct {
	TotalCustomer int `json:"total_customer"`
}
