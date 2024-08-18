package service

import (
	"billing-engine/internal/payment/model"
	"billing-engine/internal/payment/repository"
	apperror "billing-engine/pkg/customerror"
	"billing-engine/pkg/enum"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/producer"
	"context"
	"github.com/google/uuid"
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
	i.log.WithField("payload", payload).Info("[ProcessPayment] processing payment")

	isExist, err := i.repo.IsCustomerHasLoan(ctx, payload.CustomerID, payload.LoanID)
	if err != nil {
		i.log.WithField("error", err).Error("[ProcessPayment] failed to check customer has loan")
		return model.ProcessPaymentResponse{}, err
	}

	if !isExist {
		return model.ProcessPaymentResponse{}, apperror.New(apperror.NotFound, "customer has no loan")
	}

	isExist, err = i.repo.IsLoanScheduleExist(ctx, payload.LoanID, payload.ScheduleID)
	if err != nil {
		i.log.WithField("error", err).Error("[ProcessPayment] failed to check loan schedule exist")
		return model.ProcessPaymentResponse{}, err
	}

	if !isExist {
		return model.ProcessPaymentResponse{}, apperror.New(apperror.NotFound, "loan schedule does not exist")
	}

	payment, err := i.repo.UpdatePaymentScheduleStatus(ctx, payload.LoanID, payload.ScheduleID, enum.PaymentStatusPaid)
	if err != nil {
		i.log.WithField("error", err).Error("[ProcessPayment] failed to update payment schedule status")
		return model.ProcessPaymentResponse{}, err
	}

	paymentEvent := model.PaymentEventPayload{
		LoanID:        payload.LoanID,
		ScheduleID:    payload.ScheduleID,
		PaymentID:     payment.PaymentID,
		PaymentStatus: payment.PaymentStatus,
		PaymentDate:   payment.PaymentDate,
	}

	producerMessage := producer.Message{
		EventID:   uuid.New().String(),
		EventName: producer.EVENT_NAME_PAYMENT_PAID,
		Data:      paymentEvent,
	}

	err = i.producer.SendMessage(ctx, producerMessage)
	if err != nil {
		i.log.WithField("error", err).Error("[ProcessPayment] failed to send message to producer")
		return model.ProcessPaymentResponse{}, err
	}

	i.log.WithField("payload", payload).Info("[ProcessPayment] payment processed")
	return model.ProcessPaymentResponse{
		AmountPaid:    payment.AmountPaid,
		PaymentID:     payment.PaymentID,
		PaymentStatus: payment.PaymentStatus,
		PaymentDate:   payment.PaymentDate,
	}, nil

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
