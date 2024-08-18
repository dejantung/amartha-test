package repository

import (
	"billing-engine/internal/payment/domain"
	"billing-engine/pkg/enum"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=../mocks/mock_payment_repository.go -package=mocks billing-engine/internal/payment/repository PaymentRepositoryProvider
type PaymentRepositoryProvider interface {
	IsLoanScheduleExist(ctx context.Context, loanID uuid.UUID, scheduleID uuid.UUID) (bool, error)
	IsCustomerHasLoan(ctx context.Context, customerID uuid.UUID, loanID uuid.UUID) (bool, error)
	UpdatePaymentScheduleStatus(ctx context.Context, loanID uuid.UUID, scheduleID uuid.UUID, status enum.PaymentStatus) (*domain.Payment, error)

	CreateLoan(ctx context.Context, loan domain.Loan) error
}
