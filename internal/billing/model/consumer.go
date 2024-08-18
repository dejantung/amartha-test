package model

import "github.com/google/uuid"

type PaymentEventPayload struct {
	CustomerID uuid.UUID `json:"customer_id"`
	ScheduleID uuid.UUID `json:"schedule_id"`
	LoanID     uuid.UUID `json:"loan_id"`
}
