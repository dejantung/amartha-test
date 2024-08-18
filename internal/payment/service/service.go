package service

import (
	"billing-engine/internal/payment/model"
	"billing-engine/internal/payment/repository"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/producer"
	"context"
)

type PaymentServiceProvider interface {
	ProcessPayment(ctx context.Context, payload model.ProcessPaymentPayload) (model.ProcessPaymentResponse, error)
	ProcessLoanEvent(ctx context.Context, payloads model.LoanCreatedPayload) error
}

type impl struct {
	repo     repository.PaymentRepositoryProvider
	producer producer.ProducerProvider
	log      logger.Logger
}

func (i impl) ProcessPayment(ctx context.Context, payload model.ProcessPaymentPayload) (model.ProcessPaymentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i impl) ProcessLoanEvent(ctx context.Context, payloads model.LoanCreatedPayload) error {
	//TODO implement me
	panic("implement me")
}

func NewPaymentService(repo repository.PaymentRepositoryProvider,
	producer producer.ProducerProvider, log logger.Logger) PaymentServiceProvider {
	return &impl{
		repo:     repo,
		log:      log,
		producer: producer,
	}
}
