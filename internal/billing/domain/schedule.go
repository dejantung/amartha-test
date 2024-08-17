package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	ScheduleID     uuid.UUID `json:"schedule_id" gorm:"type:uuid;primaryKey"`
	LoanID         uuid.UUID `json:"loan_id" gorm:"type:uuid"`
	PaymentNo      int       `json:"payment_no"`
	PaymentDueDate time.Time `json:"payment_due_date"`
	PaymentAmount  float64   `json:"payment_amount"`
	PaymentStatus  string    `json:"payment_status"`
	IsMissPayment  bool      `json:"is_miss_payment"`
	AuditLog

	Loan Loan `json:"loan" gorm:"foreignKey:LoanID"`
}

func (schedule *Schedule) BeforeCreate(tx *gorm.DB) (err error) {
	schedule.ScheduleID = uuid.New()
	return schedule.AuditLog.BeforeCreate(tx)
}
