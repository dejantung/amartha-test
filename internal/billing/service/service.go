package service

import (
	"billing-engine/internal/billing/constant"
	"billing-engine/internal/billing/domain"
	"billing-engine/internal/billing/model"
	"billing-engine/internal/billing/repository"
	apperror "billing-engine/pkg/customerror"
	"billing-engine/pkg/enum"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/producer"
	"context"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"math"
	"time"
)

type BillingServiceProvider interface {
	CreateLoan(ctx context.Context, payload model.CreateLoanPayload) (*model.CreateLoanResponse, error)
	GetPaymentSchedule(ctx context.Context, request model.GetSchedulePayload) (*model.GetScheduleResponse, error)
	IsCustomerDelinquency(ctx context.Context, customerID uuid.UUID) (*model.IsDelinquentResponse, error)
	GetOutstandingBalance(ctx context.Context, customerID uuid.UUID) (*model.GetOutstandingBalanceResponse, error)
	ProcessMessage(ctx context.Context, payload []byte) error

	// CreateCustomer NOTE: this method is out of context, so I will just merge it in the billing service
	CreateCustomer(ctx context.Context, payload model.CreateCustomerPayload) (*model.GetCustomerResponse, error)
	// GetCustomer NOTE: this method is out of context, so I will just merge it in the billing service
	GetCustomer(ctx context.Context) ([]domain.Customer, error)
}

type BillingService struct {
	repo     repository.BillingRepositoryProvider
	log      logger.Logger
	cache    repository.BillingCacheProvider
	producer producer.ProducerProvider
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
			WithField("error", err.Error()).Info("[CreateLoan] Unexpected error when creating loan")
		return nil, err
	}

	b.log.WithField("customer_id", payload.CustomerID).
		WithField("loan", loan).Info("[CreateLoan] loan created successfully")

	newEventID := uuid.New().String()
	producerMessage := producer.Message{
		EventID:   newEventID,
		EventName: producer.EVENT_NAME_LOAN_CREATED,
		Data:      newLoan,
	}

	b.log.WithField("customer_id", payload.CustomerID).
		WithField("producer_payload", producerMessage).Info("[CreateLoan] sending message to producer")
	err = b.producer.SendMessage(ctx, producerMessage)
	if err != nil {
		b.log.WithField("customer_id", payload.CustomerID).
			WithField("error", err.Error()).Error("[CreateLoan][SendMessage] failed to send message to producer")
		return nil, err
	}

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

func (b BillingService) GetPaymentSchedule(ctx context.Context, request model.GetSchedulePayload) (*model.GetScheduleResponse, error) {
	b.log.WithField("loan_id", request.LoanID).
		WithField("customer_id", request.CustomerID).Info("[GetPaymentSchedule] getting payment schedule for loan")

	loan, err := b.repo.GetLoanByIDAndCustomerID(ctx, request.LoanID, request.CustomerID)
	if err != nil {
		b.log.WithField("loan_id", request.LoanID).
			WithField("customer_id", request.CustomerID).
			WithField("error", err.Error()).Error("[GetPaymentSchedule] Unexpected error when getting loan")
		return nil, err
	}

	if loan == nil {
		b.log.WithField("loan_id", request.LoanID).
			WithField("customer_id", request.CustomerID).Info("[GetPaymentSchedule] loan not found")
		return nil, apperror.New(apperror.NotFound, "loan not found")
	}

	b.log.WithField("loan_id", request.LoanID).
		WithField("customer_id", request.CustomerID).Info("[GetPaymentSchedule] loan found")

	return &model.GetScheduleResponse{
		Schedules: b.MapScheduleResponse(loan.Schedules),
	}, nil
}

func (b BillingService) IsCustomerDelinquency(ctx context.Context, customerID uuid.UUID) (*model.IsDelinquentResponse, error) {
	b.log.WithField("customer_id", customerID).Info("[IsCustomerDelinquency] checking customer in cache")
	resp := &model.IsDelinquentResponse{}

	cacheKey := fmt.Sprintf(constant.CACHE_KEY_DELIQUENCY, customerID)
	cacheData, err := b.cache.Get(ctx, cacheKey)
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[IsCustomerDelinquency - Cache Get] Unexpected error when getting cache")
		return nil, err
	}

	if cacheData != nil {
		b.log.WithField("customer_id", customerID).Info("[IsCustomerDelinquency] customer found in cache")
		err = json.Unmarshal([]byte(cacheData.(string)), resp)
		if err != nil {
			b.log.WithField("customer_id", customerID).
				WithField("error", err.Error()).Error("[IsCustomerDelinquency - Cache Unmarshal] Unexpected error when unmarshal cache")
			return nil, err
		}

		return resp, nil
	}

	defer func() {
		if err == nil {
			cacheErr := b.cache.Set(ctx, cacheKey, resp)
			if cacheErr != nil {
				b.log.WithField("customer_id", customerID).
					WithField("error", cacheErr.Error()).Error("[IsCustomerDelinquency - Cache Set] Unexpected error when setting cache")
				err = cacheErr
			}
		}
	}()

	customer, err := b.repo.GetCustomerByID(ctx, customerID)
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[GetCustomerByID] Unexpected error when getting customer")
		return nil, err
	}

	if customer == nil {
		b.log.WithField("customer_id", customerID).Info("[IsCustomerDelinquency] customer not found")
		return nil, apperror.New(apperror.NotFound, "customer not found")
	}

	latestLoan, err := b.repo.LastActiveLoan(ctx, customerID)
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[GetLatestActiveLoan] Unexpected error when getting loan")
		return nil, err
	}

	// we only get the unpaid and miss payment until now
	loanSchedule, err := b.repo.GetUnpaidAndMissPaymentUntil(ctx, latestLoan.LoanID, time.Now())
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[GetLatestActiveLoan] Unexpected error when getting loan")
		return nil, err
	}

	if loanSchedule == nil || len(loanSchedule) < 2 {
		return resp, nil
	}

	// since we only get the unpaid and miss payment until now, we can assume that the loan schedule is sorted,
	// so we can just check the difference between the payment number
	// if the difference is 1 then the customer is delinquent
	for i := 0; i < len(loanSchedule)-1; i++ {
		if loanSchedule[i+1].PaymentNo-loanSchedule[i].PaymentNo == 1 {
			resp.IsDelinquent = true
			break
		}
	}

	return resp, err
}

func (b BillingService) GetOutstandingBalance(ctx context.Context, customerID uuid.UUID) (*model.GetOutstandingBalanceResponse, error) {
	b.log.WithField("customer_id", customerID).Info("[GetOutstandingBalance] getting outstanding balance for customer")
	var resp model.GetOutstandingBalanceResponse

	cacheKey := fmt.Sprintf(constant.CACHE_KEY_OUTSTANDING, customerID)
	cacheData, err := b.cache.Get(ctx, cacheKey)
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[GetOutstandingBalance - Cache Get] Unexpected error when getting cache")
		return nil, err
	}

	if cacheData != nil {
		err = json.Unmarshal([]byte(cacheData.(string)), &resp)
		if err != nil {
			b.log.WithField("customer_id", customerID).
				WithField("error", err.Error()).Error("[GetOutstandingBalance - Cache Unmarshal] Unexpected error when unmarshal cache")
			return nil, err
		}

		b.log.WithField("customer_id", customerID).Info("[GetOutstandingBalance] customer found in cache")
		return &resp, nil
	}

	customer, err := b.repo.GetCustomerByID(ctx, customerID)
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[GetCustomerByID] Unexpected error when getting customer")
		return nil, err
	}

	if customer == nil {
		b.log.WithField("customer_id", customerID).Info("[GetOutstandingBalance] customer not found")
		return nil, apperror.New(apperror.NotFound, "customer not found")
	}

	lastActiveLoan, err := b.repo.LastActiveLoan(ctx, customerID)
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[GetLatestActiveLoan] Unexpected error when getting loan")
		return nil, err
	}

	totalOutstandingBalance, err := b.repo.GetTotalUnpaidPaymentOnActiveLoan(ctx, lastActiveLoan.LoanID)
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[GetTotalOutstandingBalance] Unexpected error when getting total outstanding balance")
		return nil, err
	}

	resp.OutstandingBalance = totalOutstandingBalance
	err = b.cache.Set(ctx, cacheKey, &resp)
	if err != nil {
		b.log.WithField("customer_id", customerID).
			WithField("error", err.Error()).Error("[GetOutstandingBalance - Cache Set] Unexpected error when setting cache")
		return nil, err
	}

	return &resp, nil
}

func (b BillingService) CreateCustomer(ctx context.Context, payload model.CreateCustomerPayload) (*model.GetCustomerResponse, error) {
	var customers []domain.Customer
	for i := 0; i < payload.TotalCustomer; i++ {
		customer := domain.Customer{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email:     gofakeit.Email(),
		}

		customers = append(customers, customer)
	}

	err := b.repo.CreateCustomer(ctx, customers)
	if err != nil {
		b.log.WithField("error", err.Error()).Error("[CreateCustomer] Unexpected error when creating customer")
		return nil, err
	}

	return nil, nil
}

func (b BillingService) GetCustomer(ctx context.Context) ([]domain.Customer, error) {
	customers, err := b.repo.GetCustomer(ctx)
	if err != nil {
		b.log.WithField("error", err.Error()).Error("[GetCustomer] Unexpected error when getting customer")
		return nil, err
	}

	return customers, nil
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

func (b BillingService) UpdatePayment(ctx context.Context, payload model.PaymentEventPayload) error {
	b.log.WithField("schedule_id", payload.ScheduleID).Info("[UpdatePayment] updating payment schedule")

	loan, err := b.repo.GetLoanByScheduleID(ctx, payload.ScheduleID)
	if err != nil {
		b.log.WithField("schedule_id", payload.ScheduleID).
			WithField("error", err.Error()).Error("[UpdatePayment] Unexpected error when getting schedule")
		return err
	}

	if loan == nil {
		b.log.WithField("schedule_id", payload.ScheduleID).Error("[UpdatePayment] loan not found")
		return apperror.New(apperror.NotFound, "loan not found")
	}

	schedule, err := b.repo.GetScheduleByID(ctx, payload.ScheduleID)
	if err != nil {
		b.log.WithField("schedule_id", payload.ScheduleID).
			WithField("error", err.Error()).Error("[UpdatePayment] Unexpected error when getting schedule")
		return err
	}

	if schedule == nil {
		b.log.WithField("schedule_id", payload.ScheduleID).Error("[UpdatePayment] schedule not found")
		return apperror.New(apperror.NotFound, "schedule not found")
	}

	if schedule.PaymentStatus == enum.PaymentStatusPaid {
		return nil
	}

	schedule.PaymentStatus = enum.PaymentStatusPaid
	err = b.repo.UpdateSchedulePayment(ctx, schedule)
	if err != nil {
		b.log.WithField("schedule_id", payload.ScheduleID).
			WithField("error", err.Error()).Error("[UpdatePayment] Unexpected error when updating schedule")
		return err
	}

	err = b.flushCache(ctx, loan.CustomerID)
	if err != nil {
		b.log.WithField("schedule_id", payload.ScheduleID).
			WithField("error", err.Error()).Error("[UpdatePayment] failed to flush cache")
		return err
	}

	b.log.WithField("schedule_id", payload.ScheduleID).Info("[UpdatePayment] schedule updated successfully")
	return nil
}

func (b BillingService) flushCache(ctx context.Context, customerID uuid.UUID) error {
	cacheKeyOutstanding := fmt.Sprintf(constant.CACHE_KEY_OUTSTANDING, customerID)
	currentValue, err := b.cache.Get(ctx, cacheKeyOutstanding)
	if err != nil {
		b.log.WithField("error", err.Error()).Error("[removeCache] failed to get cache")
		return err
	}

	if currentValue != nil {
		b.log.WithField("key", cacheKeyOutstanding).Info("[removeCache] key found in cache")
		err = b.cache.Delete(ctx, cacheKeyOutstanding)
		if err != nil {
			b.log.WithField("error", err.Error()).Error("[removeCache] failed to delete key from cache")
			return err
		}
	}

	cacheKeyDelinquency := fmt.Sprintf(constant.CACHE_KEY_DELIQUENCY, customerID)
	currentValue, err = b.cache.Get(ctx, cacheKeyDelinquency)
	if err != nil {
		b.log.WithField("error", err.Error()).Error("[removeCache] failed to get cache")
		return err
	}

	if currentValue != nil {
		b.log.WithField("key", cacheKeyDelinquency).Info("[removeCache] key found in cache")
		err = b.cache.Delete(ctx, cacheKeyDelinquency)
		if err != nil {
			b.log.WithField("error", err.Error()).Error("[removeCache] failed to delete key from cache")
			return err
		}
	}

	return nil
}

func (b BillingService) ProcessMessage(ctx context.Context, payload []byte) error {
	b.log.WithField("payload", string(payload)).Info("[ProcessMessage] processing message")

	var message producer.Message
	err := json.Unmarshal(payload, &message)
	if err != nil {
		b.log.WithField("error", err).Error("[ProcessMessage] failed to unmarshal message payload")
		return err
	}

	switch message.EventName {
	case producer.EVENT_NAME_PAYMENT_PAID:
		var parseData model.PaymentEventPayload

		dataByte, err := json.Marshal(message.Data)
		if err != nil {
			b.log.WithField("error", err).Error("[ProcessMessage] failed to marshal message.Data")
			return err
		}
		err = json.Unmarshal(dataByte, &parseData)
		if err != nil {
			b.log.WithField("error", err).Error("[ProcessMessage] failed to assert message.Data to model")
			return err
		}

		err = b.UpdatePayment(ctx, parseData)
		if err != nil {
			b.log.WithField("error", err).Error("[ProcessMessage] failed to process loan event")
			return err
		}
	default:
		b.log.WithField("event_name", message.EventName).
			WithField("payload", message).Error("[ProcessMessage] unknown event name")
	}

	b.log.WithField("payload", string(payload)).Info("[ProcessMessage] message processed")
	return nil
}

func NewBillingService(repo repository.BillingRepositoryProvider,
	cache repository.BillingCacheProvider, producer producer.ProducerProvider, log logger.Logger) *BillingService {
	return &BillingService{
		repo:     repo,
		log:      log,
		cache:    cache,
		producer: producer,
	}
}
