package repository

import (
	"billing-engine/internal/billing/domain"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=../mocks/mock_billing_repository.go -package=mocks billing-engine/internal/billing/repository BillingRepositoryProvider
type BillingRepositoryProvider interface {
	CreateLoan(ctx context.Context, request domain.Loan) (*domain.Loan, error)
	GetSchedule(ctx context.Context, loanID, customerID uuid.UUID) ([]domain.Schedule, error)
	GetAllSchedule(ctx context.Context, customerID uuid.UUID) ([]domain.Schedule, error)
	GetLatestLoan(ctx context.Context, customerID uuid.UUID) (*domain.Loan, error)

	// CreateCustomer NOTE: this method is out of context, so I will just merge it in the billing service
	CreateCustomer(ctx context.Context, request domain.Customer) (*domain.Customer, error)
	// GetCustomer NOTE: this method is out of context, so I will just merge it in the billing service
	GetCustomer(ctx context.Context) ([]domain.Customer, error)
}
