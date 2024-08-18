package service_test

import (
	"billing-engine/internal/payment/domain"
	"billing-engine/internal/payment/mocks"
	"billing-engine/internal/payment/model"
	"billing-engine/pkg/logger"
	pkgMock "billing-engine/pkg/mocks"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"billing-engine/internal/payment/service"
)

var someErr = errors.New("some error")

var _ = Describe("Service", func() {
	var (
		svc      service.PaymentServiceProvider
		mockCtrl *gomock.Controller
		repo     *mocks.MockPaymentRepositoryProvider
		log      logger.Logger
		producer *pkgMock.MockProducerProvider
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		log = logger.NewZeroLogger("tests")
		repo = mocks.NewMockPaymentRepositoryProvider(mockCtrl)
		producer = pkgMock.NewMockProducerProvider(mockCtrl)
		svc = service.NewPaymentService(repo, producer, log)
	})

	Describe("ProcessPayment", func() {
		payload := model.ProcessPaymentPayload{}
		mockPayment := &domain.Payment{}

		Describe("Positive Case", func() {
			It("when payment is successful", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().IsLoanScheduleExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().UpdatePaymentScheduleStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockPayment, nil)
				producer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).To(BeNil())
			})
		})

		Describe("Negative Case", func() {
			It("when customer has no loan", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).ToNot(BeNil())
			})

			It("when loan schedule does not exist", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().IsLoanScheduleExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).ToNot(BeNil())
			})

			It("when update payment schedule status failed", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().IsLoanScheduleExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().UpdatePaymentScheduleStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, someErr)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).To(HaveOccurred())
			})

			It("when sending message to kafka failed", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().IsLoanScheduleExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().UpdatePaymentScheduleStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockPayment, nil)
				producer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).To(BeNil())
			})

			It("when error getting customer loan", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, someErr)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).To(HaveOccurred())
			})

			It("when error getting loan schedule", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().IsLoanScheduleExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, someErr)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).To(HaveOccurred())
			})

			It("when error updating payment schedule status", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().IsLoanScheduleExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().UpdatePaymentScheduleStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, someErr)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).To(HaveOccurred())
			})

			It("when sending message to producer failed", func() {
				repo.EXPECT().IsCustomerHasLoan(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().IsLoanScheduleExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				repo.EXPECT().UpdatePaymentScheduleStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockPayment, nil)
				producer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(someErr)

				_, err := svc.ProcessPayment(nil, payload)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ProcessLoanEvent", func() {
		payload := model.LoanCreatedPayload{}
		mockLoan := domain.Loan{}

		Describe("Positive Case", func() {
			It("when loan event is successfully processed", func() {
				repo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(mockLoan, nil)

				err := svc.ProcessLoanEvent(nil, payload)
				Expect(err).To(BeNil())
			})
		})

		Describe("Negative Case", func() {
			It("when error processing loan event", func() {
				repo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(mockLoan, someErr)

				err := svc.ProcessLoanEvent(nil, payload)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
