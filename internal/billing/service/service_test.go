package service_test

import (
	"billing-engine/internal/billing/domain"
	"billing-engine/internal/billing/mocks"
	"billing-engine/internal/billing/model"
	"billing-engine/internal/billing/service"
	apperror "billing-engine/pkg/customerror"
	"billing-engine/pkg/enum"
	"billing-engine/pkg/logger"
	"context"
	"errors"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"time"
)

var someErr = errors.New("some error")
var ctx = context.Background()
var randUUID = uuid.New()
var timeNow = time.Now()

var _ = Describe("Service", func() {
	var (
		mockCtrl     *gomock.Controller
		svc          service.BillingServiceProvider
		repo         *mocks.MockBillingRepositoryProvider
		log          logger.Logger
		mockLoan     domain.Loan
		mockSchedule []domain.Schedule
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repo = mocks.NewMockBillingRepositoryProvider(mockCtrl)
		log = logger.NewZeroLogger("test")
		svc = service.NewBillingService(repo, log)

		mockSchedule = []domain.Schedule{
			{
				ScheduleID:     uuid.New(),
				LoanID:         randUUID,
				PaymentNo:      1,
				PaymentDueDate: timeNow.AddDate(0, 1, 0),
				PaymentAmount:  1,
				PaymentStatus:  enum.PaymentStatusPending,
				IsMissPayment:  false,
			},
			{
				ScheduleID:     uuid.New(),
				LoanID:         randUUID,
				PaymentNo:      2,
				PaymentDueDate: timeNow.AddDate(0, 2, 0),
				PaymentAmount:  1,
				PaymentStatus:  enum.PaymentStatusPending,
				IsMissPayment:  false,
			},
			{
				ScheduleID:     uuid.New(),
				LoanID:         randUUID,
				PaymentNo:      3,
				PaymentDueDate: timeNow.AddDate(0, 3, 0),
				PaymentAmount:  1,
				PaymentStatus:  enum.PaymentStatusPending,
				IsMissPayment:  false,
			},
			{
				ScheduleID:     uuid.New(),
				LoanID:         randUUID,
				PaymentNo:      4,
				PaymentDueDate: timeNow.AddDate(0, 4, 0),
				PaymentAmount:  1,
				PaymentStatus:  enum.PaymentStatusPending,
				IsMissPayment:  false,
			},
			{
				ScheduleID:     uuid.New(),
				LoanID:         randUUID,
				PaymentNo:      5,
				PaymentDueDate: timeNow.AddDate(0, 5, 0),
				PaymentAmount:  1,
				PaymentStatus:  enum.PaymentStatusPending,
				IsMissPayment:  false,
			},
		}

		mockLoan = domain.Loan{
			LoanID:             randUUID,
			CustomerID:         uuid.New(),
			PrincipalAmount:    5000000,
			InterestRate:       0.1,
			StartDate:          timeNow,
			EndDate:            timeNow.AddDate(0, 5, 0),
			OutstandingBalance: 0,
			IsDelinquent:       false,
			Customer: domain.Customer{
				CustomerID: uuid.New(),
			},
			AuditLog: domain.AuditLog{},
		}
	})

	Describe("CreateLoan", func() {
		payload := model.CreateLoanPayload{
			CustomerID: randUUID,
			LoanAmount: 5000000,
		}

		Describe("Positive case", func() {
			It("should return correct loan response", func() {
				repo.EXPECT().GetCustomerByID(ctx, payload.CustomerID).Return(&domain.Customer{}, nil)
				mockLoan.Schedules = mockSchedule
				repo.EXPECT().CreateLoan(ctx, gomock.Any()).Return(&mockLoan, nil)

				response, err := svc.CreateLoan(ctx, payload)
				Expect(err).To(BeNil())
				Expect(len(response.Schedules)).To(Equal(len(mockSchedule)))
				Expect(response.LoanAmount).To(Equal(float64(5500000)))

				for i, val := range response.Schedules {
					Expect(val.PaymentNo).To(Equal(mockSchedule[i].PaymentNo))
					Expect(val.PaymentDueDate).To(Equal(mockSchedule[i].PaymentDueDate.Format("2006-01-02")))
					Expect(val.PaymentAmount).To(Equal(mockSchedule[i].PaymentAmount))
					Expect(val.PaymentStatus).To(Equal(mockSchedule[i].PaymentStatus))
					Expect(val.IsMissPayment).To(Equal(mockSchedule[i].IsMissPayment))
				}
			})
		})

		Describe("Negative case", func() {
			It("when customer not found", func() {
				repo.EXPECT().GetCustomerByID(ctx, payload.CustomerID).Return(nil, nil)
				_, err := svc.CreateLoan(ctx, payload)

				var errs *apperror.CustomError
				ok := errors.As(err, &errs)
				Expect(ok).To(BeTrue())
				Expect(errs.Cause).To(Equal(apperror.NotFound))
			})

			It("when error on get customer by id", func() {
				repo.EXPECT().GetCustomerByID(ctx, payload.CustomerID).Return(nil, someErr)
				_, err := svc.CreateLoan(ctx, payload)
				Expect(err).To(Equal(someErr))
			})

			It("when error on create loan", func() {
				repo.EXPECT().GetCustomerByID(ctx, payload.CustomerID).Return(&domain.Customer{}, nil)
				repo.EXPECT().CreateLoan(ctx, gomock.Any()).Return(nil, someErr)
				_, err := svc.CreateLoan(ctx, payload)
				Expect(err).To(Equal(someErr))
			})
		})
	})
})
