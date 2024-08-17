package service

import (
	"billing-engine/internal/billing/model"
	"billing-engine/internal/billing/repository"
	"billing-engine/pkg/logger"
	"context"
)

type BillingServiceProvider interface {
	CreateLoan(ctx context.Context, payload model.CreateLoanPayload) (*model.CreateLoanResponse, error)
	GetPaymentSchedule(ctx context.Context, request model.GetScheduleResponse) (*model.GetScheduleResponse, error)
	IsCustomerDelinquency(ctx context.Context, payload model.IsDelinquentPayload) (*model.IsDelinquentResponse, error)
	GetOutstandingBalance(ctx context.Context, payload model.IsDelinquentPayload) (float64, error)

	// CreateCustomer NOTE: this method is out of context, so I will just merge it in the billing service
	CreateCustomer(ctx context.Context, payload model.CreateCustomerPayload) (*model.GetCustomerResponse, error)
	// GetCustomer NOTE: this method is out of context, so I will just merge it in the billing service
	GetCustomer(ctx context.Context) (*model.GetCustomerResponse, error)
}

type BillingService struct {
	repo repository.BillingRepositoryProvider
	log  logger.Logger
}

func (b BillingService) CreateLoan(ctx context.Context, payload model.CreateLoanPayload) (*model.CreateLoanResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b BillingService) GetPaymentSchedule(ctx context.Context, request model.GetScheduleResponse) (*model.GetScheduleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b BillingService) IsCustomerDelinquency(ctx context.Context, payload model.IsDelinquentPayload) (*model.IsDelinquentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b BillingService) GetOutstandingBalance(ctx context.Context, payload model.IsDelinquentPayload) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (b BillingService) CreateCustomer(ctx context.Context, payload model.CreateCustomerPayload) (*model.GetCustomerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b BillingService) GetCustomer(ctx context.Context) (*model.GetCustomerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewBillingService(repo repository.BillingRepositoryProvider, log logger.Logger) BillingServiceProvider {
	return &BillingService{
		repo: repo,
		log:  log,
	}
}
