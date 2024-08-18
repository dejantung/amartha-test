package repository

import (
	"billing-engine/internal/billing/domain"
	"billing-engine/pkg/enum"
	"billing-engine/pkg/logger"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_billing_repository.go -package=mocks billing-engine/internal/billing/repository BillingRepositoryProvider
type BillingRepositoryProvider interface {
	CreateLoan(ctx context.Context, request domain.Loan) (*domain.Loan, error)
	GetSchedule(ctx context.Context, loanID, customerID uuid.UUID) ([]domain.Schedule, error)
	GetUnpaidAndMissPaymentUntil(ctx context.Context, customerID uuid.UUID, date time.Time) ([]domain.Schedule, error)
	GetLoanByIDAndCustomerID(ctx context.Context, loanID, customerID uuid.UUID) (*domain.Loan, error)
	GetTotalUnpaidPaymentOnActiveLoan(ctx context.Context, loanId uuid.UUID) (float64, error)
	LastActiveLoan(ctx context.Context, customerID uuid.UUID) (*domain.Loan, error)

	// CreateCustomer NOTE: this method is out of context, so I will just merge it in the billing service
	CreateCustomer(ctx context.Context, request []domain.Customer) error
	// GetCustomer NOTE: this method is out of context, so I will just merge it in the billing service
	GetCustomer(ctx context.Context) ([]domain.Customer, error)
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*domain.Customer, error)
}

type repo struct {
	db  *gorm.DB
	log logger.Logger
}

func (r repo) CreateLoan(ctx context.Context, request domain.Loan) (*domain.Loan, error) {
	err := r.db.WithContext(ctx).Create(&request).Error
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r repo) GetSchedule(ctx context.Context, loanID, customerID uuid.UUID) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.WithContext(ctx).Where("loan_id = ? AND customer_id = ?", loanID, customerID).Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r repo) GetUnpaidAndMissPaymentUntil(ctx context.Context, customerID uuid.UUID, date time.Time) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.WithContext(ctx).Where("customer_id = ? AND due_date < ? AND status = ?", customerID, date, enum.PaymentStatusPending).Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r repo) GetLoanByIDAndCustomerID(ctx context.Context, loanID, customerID uuid.UUID) (*domain.Loan, error) {
	var loan domain.Loan
	err := r.db.WithContext(ctx).Where("loan_id = ? AND customer_id = ?", loanID, customerID).First(&loan).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r repo) LastActiveLoan(ctx context.Context, customerID uuid.UUID) (*domain.Loan, error) {
	var loan domain.Loan
	err := r.db.WithContext(ctx).Where("customer_id = ? AND is_finish = ?", customerID, false).First(&loan).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r repo) GetTotalUnpaidPaymentOnActiveLoan(ctx context.Context, loanId uuid.UUID) (float64, error) {
	var totalUnpaid float64
	err := r.db.WithContext(ctx).Model(&domain.Schedule{}).
		Select("SUM(payment_amount)").
		Where("loan_id = ? AND payment_status = ?", loanId, enum.PaymentStatusPending).
		Row().
		Scan(&totalUnpaid)
	if err != nil {
		return 0, err
	}

	return totalUnpaid, nil
}

func (r repo) CreateCustomer(ctx context.Context, request []domain.Customer) error {
	return r.db.WithContext(ctx).Create(&request).Error
}

func (r repo) GetCustomer(ctx context.Context) ([]domain.Customer, error) {
	var result []domain.Customer

	err := r.db.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r repo) GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).First(&customer).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &customer, nil
}

func NewBillingRepositoryProvider(db *gorm.DB, log logger.Logger) BillingRepositoryProvider {
	return &repo{
		db:  db,
		log: log,
	}
}
