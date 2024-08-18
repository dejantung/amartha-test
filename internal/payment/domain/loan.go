package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Loan struct {
	Base
	LoanID     uuid.UUID `json:"loan_id" gorm:"type:uuid;primaryKey"`
	CustomerID uuid.UUID `json:"customer_id" gorm:"type:uuid"`

	PaymentSchedules []PaymentSchedule `json:"payment_schedules" gorm:"foreignKey:LoanID"`
}

func (loan *Loan) BeforeCreate(tx *gorm.DB) (err error) {
	return loan.Base.BeforeCreate(tx)
}
