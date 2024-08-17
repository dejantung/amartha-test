package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Loan struct {
	LoanID             uuid.UUID `json:"loan_id"`
	CustomerID         uuid.UUID `json:"customer_id"`
	PrincipalAmount    float64   `json:"principal_amount"`
	InterestRate       float64   `json:"interest_rate"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	OutstandingBalance float64   `json:"outstanding_balance"`
	IsDelinquent       bool      `json:"is_delinquent"`

	Customer  Customer   `json:"customer" gorm:"foreignKey:CustomerID"`
	Schedules []Schedule `json:"schedules" gorm:"foreignKey:LoanID"`
	AuditLog
}

func (loan *Loan) BeforeCreate(tx *gorm.DB) (err error) {
	return loan.AuditLog.BeforeCreate(tx)
}
