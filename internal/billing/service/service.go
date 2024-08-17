package service

import (
	"billing-engine/internal/billing/constant"
	"billing-engine/internal/billing/domain"
	"billing-engine/internal/billing/model"
	"billing-engine/internal/billing/repository"
	apperror "billing-engine/pkg/customerror"
	"billing-engine/pkg/enum"
	"billing-engine/pkg/logger"
	"context"
	"math"
	"time"
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
	b.log.WithField("customer_id", payload.CustomerID).Info("[CreateLoan] creating loan for customer")
	customer, err := b.repo.GetCustomerByID(ctx, payload.CustomerID)
	if err != nil {
		b.log.WithField("customer_id", payload.CustomerID).
			WithField("error", err.Error()).Error("[CreateLoan] Unexpected error when getting customer")
		return nil, err
	}

	if customer == nil {
		b.log.WithField("customer_id", payload.CustomerID).Error("[CreateLoan] customer not found")
		return nil, apperror.New(apperror.NotFound, "customer not found")
	}

	var loanSchema []domain.Schedule

	loan := domain.Loan{
		CustomerID:      payload.CustomerID,
		PrincipalAmount: payload.LoanAmount,
		InterestRate:    constant.INTEREST_RATE,
		StartDate:       time.Now(),
		EndDate:         time.Now().AddDate(0, 5, 0),
	}

	totalLoan, loanSchema := b.paymentSchemaMaker(loan)
	loan.Schedules = loanSchema

	b.log.WithField("customer_id", payload.CustomerID).Info("[CreateLoan] creating loan for customer")
	newLoan, err := b.repo.CreateLoan(ctx, loan)
	if err != nil {
		b.log.WithField("customer_id", payload.CustomerID).
			WithField("error", err.Error()).Error("[CreateLoan] Unexpected error when creating loan")
		return nil, err
	}

	//TODO: add kafka producer here
	b.log.WithField("customer_id", payload.CustomerID).
		WithField("loan", loan).Info("[CreateLoan] loan created successfully")

	return &model.CreateLoanResponse{
		LoanID:     newLoan.LoanID,
		CustomerID: newLoan.CustomerID,
		LoanAmount: totalLoan,
		Schedules:  b.MapScheduleResponse(newLoan.Schedules),
	}, nil
}

func (b BillingService) paymentSchemaMaker(loan domain.Loan) (float64, []domain.Schedule) {
	var newSchedule []domain.Schedule
	totalLoan := math.Round(loan.PrincipalAmount + (loan.PrincipalAmount * loan.InterestRate))

	paymentPerMonth := math.Round(totalLoan / constant.MAX_PAYMENT)

	for i := 1; i <= constant.MAX_PAYMENT; i++ {
		newSchedule = append(newSchedule, domain.Schedule{
			PaymentNo:      i,
			PaymentDueDate: time.Now().AddDate(0, i, 0),
			PaymentAmount:  math.Round(paymentPerMonth),
			PaymentStatus:  enum.PaymentStatusPending,
			IsMissPayment:  false,
		})
	}

	return totalLoan, newSchedule
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

func (b BillingService) MapScheduleResponse(schedule []domain.Schedule) []model.ScheduleResponse {
	var scheduleResp []model.ScheduleResponse

	for _, val := range schedule {
		scheduleResp = append(scheduleResp, model.ScheduleResponse{
			ScheduleID:     val.ScheduleID,
			LoanID:         val.LoanID,
			PaymentNo:      val.PaymentNo,
			PaymentDueDate: val.PaymentDueDate.Format("2006-01-02"),
			PaymentAmount:  val.PaymentAmount,
			PaymentStatus:  val.PaymentStatus,
			IsMissPayment:  val.IsMissPayment,
		})
	}

	return scheduleResp
}

func NewBillingService(repo repository.BillingRepositoryProvider, log logger.Logger) *BillingService {
	return &BillingService{
		repo: repo,
		log:  log,
	}
}
