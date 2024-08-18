package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Loan struct {
	LoanID          uuid.UUID `json:"loan_id" gorm:"type:uuid;primaryKey"`
	CustomerID      uuid.UUID `json:"customer_id" gorm:"type:uuid"`
	PrincipalAmount float64   `json:"principal_amount"`
	InterestRate    float64   `json:"interest_rate"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	IsFinish        bool      `json:"is_finish"`

	Schedules []Schedule `json:"schedules" gorm:"foreignKey:LoanID;references:LoanID"`
	AuditLog
}

func (loan *Loan) BeforeCreate(tx *gorm.DB) (err error) {
	return loan.AuditLog.BeforeCreate(tx)
}
