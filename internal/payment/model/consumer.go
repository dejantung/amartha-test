package model

import (
	"billing-engine/pkg/enum"
	"github.com/google/uuid"
	"time"
)

type LoanCreatedPayload struct {
	LoanID     uuid.UUID `json:"loan_id"`
	CustomerID uuid.UUID `json:"customer_id"`

	Schedules []LoanSchedule `json:"schedules"`
}

type LoanSchedule struct {
	ScheduleID     uuid.UUID          `json:"schedule_id"`
	PaymentNo      int                `json:"payment_no"`
	PaymentDueDate time.Time          `json:"payment_due_date"`
	PaymentAmount  float64            `json:"payment_amount"`
	PaymentStatus  enum.PaymentStatus `json:"payment_status"`
}
