package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PaymentSchedule struct {
	Base
	ScheduleID     uuid.UUID `json:"schedule_id" gorm:"type:uuid;primaryKey"`
	LoanID         uuid.UUID `json:"loan_id" gorm:"type:uuid"`
	PaymentNo      int       `json:"payment_no"`
	PaymentDueDate time.Time `json:"payment_due_date"`
	PaymentAmount  float64   `json:"payment_amount"`
	PaymentStatus  string    `json:"payment_status"`

	Payment Payment `json:"payments" gorm:"foreignKey:ScheduleID"`
}

func (paymentSchedule *PaymentSchedule) BeforeCreate(tx *gorm.DB) (err error) {
	paymentSchedule.ScheduleID = uuid.New()
	return paymentSchedule.Base.BeforeCreate(tx)
}
