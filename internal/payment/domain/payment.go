package domain

import (
	"billing-engine/pkg/enum"
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	Base
	PaymentID     uuid.UUID          `json:"payment_id" gorm:"type:uuid;primaryKey"`
	LoanID        uuid.UUID          `json:"loan_id" gorm:"type:uuid"`
	ScheduleID    uuid.UUID          `json:"schedule_id" gorm:"type:uuid"`
	PaymentDate   time.Time          `json:"payment_date"`
	AmountPaid    float64            `json:"amount_paid"`
	PaymentMethod string             `json:"payment_method"`
	PaymentStatus enum.PaymentStatus `json:"payment_status"`
}
