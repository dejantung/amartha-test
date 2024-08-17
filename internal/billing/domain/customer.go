package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	CustomerID  uuid.UUID `json:"customer_id" gorm:"type:uuid;primaryKey"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email" gorm:"uniqueIndex"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Loans []Loan `json:"loans" gorm:"foreignKey:CustomerID"`
	AuditLog
}

func (customer *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	customer.CustomerID = uuid.New()
	return customer.AuditLog.BeforeCreate(tx)
}
