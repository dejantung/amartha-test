package model

import (
	"billing-engine/pkg/enum"
	"github.com/google/uuid"
	"time"
)

type ProcessPaymentPayload struct {
	Amount     float64   `json:"amount"`
	LoanID     uuid.UUID `json:"loan_id"`
	ScheduleID uuid.UUID `json:"schedule_id"`
	CustomerID uuid.UUID `json:"customer_id"`
}

type ProcessPaymentResponse struct {
	AmountPaid    float64            `json:"amount_paid"`
	PaymentID     uuid.UUID          `json:"payment_id"`
	PaymentStatus enum.PaymentStatus `json:"payment_status"`
	PaymentDate   string             `json:"payment_date"`
}

type PaymentEventPayload struct {
	LoanID        uuid.UUID          `json:"loan_id"`
	ScheduleID    uuid.UUID          `json:"schedule_id"`
	PaymentID     uuid.UUID          `json:"payment_id"`
	PaymentStatus enum.PaymentStatus `json:"payment_status"`
	PaymentDate   time.Time          `json:"payment_date"`
}
