package repository

import (
	"billing-engine/internal/payment/domain"
	"billing-engine/pkg/enum"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../mocks/mock_payment_repository.go -package=mocks billing-engine/internal/payment/repository PaymentRepositoryProvider
type PaymentRepositoryProvider interface {
	IsLoanScheduleExist(ctx context.Context, loanID uuid.UUID, scheduleID uuid.UUID) (bool, error)
	IsCustomerHasLoan(ctx context.Context, customerID uuid.UUID, loanID uuid.UUID) (bool, error)
	UpdatePaymentScheduleStatus(ctx context.Context, loanID uuid.UUID, scheduleID uuid.UUID, status enum.PaymentStatus) (*domain.PaymentSchedule, error)
	CreatePayment(ctx context.Context, payment domain.Payment) (domain.Payment, error)

	CreateLoan(ctx context.Context, loan domain.Loan) (domain.Loan, error)
}

type impl struct {
	db *gorm.DB
}

func (i impl) IsLoanScheduleExist(ctx context.Context, loanID uuid.UUID, scheduleID uuid.UUID) (bool, error) {
	var count int64
	err := i.db.WithContext(ctx).Model(&domain.PaymentSchedule{}).
		Where("loan_id = ? AND schedule_id = ?", loanID, scheduleID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (i impl) IsCustomerHasLoan(ctx context.Context, customerID uuid.UUID, loanID uuid.UUID) (bool, error) {
	var count int64
	err := i.db.WithContext(ctx).Model(&domain.Loan{}).
		Where("customer_id = ? AND loan_id = ?", customerID, loanID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (i impl) UpdatePaymentScheduleStatus(ctx context.Context, loanID uuid.UUID, scheduleID uuid.UUID, status enum.PaymentStatus) (*domain.PaymentSchedule, error) {
	var payment domain.PaymentSchedule

	// Retrieve the payment schedule record
	err := i.db.WithContext(ctx).Model(&domain.PaymentSchedule{}).
		Where("loan_id = ? AND schedule_id = ?", loanID, scheduleID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	payment.PaymentStatus = status

	err = i.db.WithContext(ctx).Save(&payment).Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (i impl) CreateLoan(ctx context.Context, loan domain.Loan) (domain.Loan, error) {
	err := i.db.WithContext(ctx).Create(&loan).Error
	if err != nil {
		return domain.Loan{}, err
	}

	return loan, nil
}

func (i impl) CreatePayment(ctx context.Context, payment domain.Payment) (domain.Payment, error) {
	err := i.db.WithContext(ctx).Create(&payment).Error
	if err != nil {
		return domain.Payment{}, err
	}

	return payment, nil
}

func NewPaymentRepository(db *gorm.DB) PaymentRepositoryProvider {
	return &impl{
		db: db,
	}
}
