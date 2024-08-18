package service

import (
	"billing-engine/internal/billing/constant"
	"billing-engine/internal/billing/domain"
	"billing-engine/internal/billing/mocks"
	"billing-engine/internal/billing/model"
	apperror "billing-engine/pkg/customerror"
	"billing-engine/pkg/enum"
	"billing-engine/pkg/logger"
	pkgMock "billing-engine/pkg/mocks"
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
		svc          *BillingService
		repo         *mocks.MockBillingRepositoryProvider
		log          logger.Logger
		mockLoan     domain.Loan
		mockSchedule []domain.Schedule
		cache        *mocks.MockBillingCacheProvider
		producer     *pkgMock.MockProducerProvider
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repo = mocks.NewMockBillingRepositoryProvider(mockCtrl)
		log = logger.NewZeroLogger("test")
		cache = mocks.NewMockBillingCacheProvider(mockCtrl)
		producer = pkgMock.NewMockProducerProvider(mockCtrl)
		svc = NewBillingService(repo, cache, producer, log)

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
			LoanID:          randUUID,
			CustomerID:      uuid.New(),
			PrincipalAmount: 5000000,
			InterestRate:    0.1,
			StartDate:       timeNow,
			EndDate:         timeNow.AddDate(0, 5, 0),
			AuditLog:        domain.AuditLog{},
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
				producer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil)

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

			It("when error on produce message", func() {
				repo.EXPECT().GetCustomerByID(ctx, payload.CustomerID).Return(&domain.Customer{}, nil)
				mockLoan.Schedules = mockSchedule
				repo.EXPECT().CreateLoan(ctx, gomock.Any()).Return(&mockLoan, nil)
				producer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(someErr)

				_, err := svc.CreateLoan(ctx, payload)
				Expect(err).To(Equal(someErr))
			})
		})
	})

	Describe("SchemaMaker", func() {
		It("should return correct total loan and schedule with even number in total amount", func() {
			totalLoan, schedules := svc.paymentSchemaMaker(mockLoan)

			Expect(totalLoan).To(Equal(float64(5500000)))
			Expect(len(schedules)).To(Equal(constant.MAX_PAYMENT))

			for i, val := range schedules {
				Expect(val.PaymentNo).To(Equal(i + 1))
				Expect(val.PaymentAmount).To(Equal(float64(110000)))
				Expect(val.PaymentStatus).To(Equal(enum.PaymentStatusPending))
				Expect(val.IsMissPayment).To(BeFalse())
			}
		})

		It("should return correct total loan and schedule with odd number in total amount", func() {
			mockLoan.PrincipalAmount = 5000001
			totalLoan, schedules := svc.paymentSchemaMaker(mockLoan)

			Expect(totalLoan).To(Equal(float64(5500001)))
			Expect(len(schedules)).To(Equal(constant.MAX_PAYMENT))

			for i, val := range schedules {
				Expect(val.PaymentNo).To(Equal(i + 1))
				Expect(val.PaymentAmount).To(Equal(float64(110000)))
				Expect(val.PaymentStatus).To(Equal(enum.PaymentStatusPending))
				Expect(val.IsMissPayment).To(BeFalse())
			}
		})

		It("should return correct total loan and schedule with more weird odd number in total amount", func() {
			mockLoan.PrincipalAmount = 1234569
			totalLoan, schedules := svc.paymentSchemaMaker(mockLoan)

			Expect(totalLoan).To(Equal(float64(1358026)))
			Expect(len(schedules)).To(Equal(constant.MAX_PAYMENT))

			for i, val := range schedules {
				Expect(val.PaymentNo).To(Equal(i + 1))
				Expect(val.PaymentAmount).To(Equal(float64(27161)))
				Expect(val.PaymentStatus).To(Equal(enum.PaymentStatusPending))
				Expect(val.IsMissPayment).To(BeFalse())
			}
		})
	})

	Describe("GetPaymentSchedule", func() {
		payload := model.GetSchedulePayload{
			LoanID:     randUUID,
			CustomerID: randUUID,
		}

		Describe("Positive case", func() {
			It("should return correct schedule response", func() {
				mockLoan.Schedules = mockSchedule
				repo.EXPECT().GetLoanByIDAndCustomerID(ctx, payload.LoanID, payload.CustomerID).Return(&mockLoan, nil)

				response, err := svc.GetPaymentSchedule(ctx, payload)
				Expect(err).To(BeNil())
				Expect(len(response.Schedules)).To(Equal(len(mockSchedule)))

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
			It("when loan not found", func() {
				repo.EXPECT().GetLoanByIDAndCustomerID(ctx, payload.LoanID, payload.CustomerID).Return(nil, nil)
				_, err := svc.GetPaymentSchedule(ctx, payload)

				var errs *apperror.CustomError
				ok := errors.As(err, &errs)
				Expect(ok).To(BeTrue())
				Expect(errs.Cause).To(Equal(apperror.NotFound))
			})

			It("when error on get loan by id and customer id", func() {
				repo.EXPECT().GetLoanByIDAndCustomerID(ctx, payload.LoanID, payload.CustomerID).Return(nil, someErr)
				_, err := svc.GetPaymentSchedule(ctx, payload)
				Expect(err).To(Equal(someErr))
			})
		})
	})

	Describe("IsDelinquent", func() {
		Describe("Positive case", func() {
			cacheRes := &model.IsDelinquentResponse{
				IsDelinquent: true,
			}

			It("should return correct response with cache", func() {
				cache.EXPECT().Get(ctx, gomock.Any()).Return(cacheRes, nil)
				response, err := svc.IsCustomerDelinquency(ctx, uuid.New())
				Expect(err).To(BeNil())
				Expect(response.IsDelinquent).To(BeTrue())
			})

			It("when customer is not delinquent", func() {
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				cache.EXPECT().Set(ctx, gomock.Any(), gomock.Any()).Return(nil)

				repo.EXPECT().GetCustomerByID(ctx, gomock.Any()).Return(&domain.Customer{}, nil)
				repo.EXPECT().GetUnpaidAndMissPaymentUntil(ctx, gomock.Any(), gomock.Any()).Return([]domain.Schedule{
					{
						PaymentNo: 1,
					},
					{
						PaymentNo: 3,
					},
					{
						PaymentNo: 5,
					},
					{
						PaymentNo: 27,
					},
				}, nil)

				response, err := svc.IsCustomerDelinquency(ctx, uuid.New())
				Expect(err).To(BeNil())
				Expect(response.IsDelinquent).To(BeFalse())
			})

			It("when customer is delinquent", func() {
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				cache.EXPECT().Set(ctx, gomock.Any(), gomock.Any()).Return(nil)

				repo.EXPECT().GetCustomerByID(ctx, gomock.Any()).Return(&domain.Customer{}, nil)
				repo.EXPECT().GetUnpaidAndMissPaymentUntil(ctx, gomock.Any(), gomock.Any()).Return([]domain.Schedule{
					{
						PaymentNo: 1,
					},
					{
						PaymentNo: 3,
					},
					{
						PaymentNo: 25,
					},
					{
						PaymentNo: 26,
					},
				}, nil)

				response, err := svc.IsCustomerDelinquency(ctx, uuid.New())
				Expect(err).To(BeNil())
				Expect(response.IsDelinquent).To(BeTrue())
			})

			It("when customer only have 1 unpaid / missing payment", func() {
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				cache.EXPECT().Set(ctx, gomock.Any(), gomock.Any()).Return(nil)

				repo.EXPECT().GetCustomerByID(ctx, gomock.Any()).Return(&domain.Customer{}, nil)
				repo.EXPECT().GetUnpaidAndMissPaymentUntil(ctx, gomock.Any(), gomock.Any()).Return([]domain.Schedule{
					{
						PaymentNo: 1,
					},
				}, nil)

				response, err := svc.IsCustomerDelinquency(ctx, uuid.New())
				Expect(err).To(BeNil())
				Expect(response.IsDelinquent).To(BeFalse())
			})
		})

		Describe("Negative case", func() {
			It("when error getting cache", func() {
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, someErr)
				_, err := svc.IsCustomerDelinquency(ctx, uuid.New())
				Expect(err).To(Equal(someErr))
			})

			It("when error getting customer by id", func() {
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				repo.EXPECT().GetCustomerByID(ctx, gomock.Any()).Return(nil, someErr)
				_, err := svc.IsCustomerDelinquency(ctx, uuid.New())
				Expect(err).To(Equal(someErr))
			})

			It("when error getting unpaid and miss payment until", func() {
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				repo.EXPECT().GetCustomerByID(ctx, gomock.Any()).Return(&domain.Customer{}, nil)
				repo.EXPECT().GetUnpaidAndMissPaymentUntil(ctx, gomock.Any(), gomock.Any()).Return(nil, someErr)
				_, err := svc.IsCustomerDelinquency(ctx, uuid.New())
				Expect(err).To(Equal(someErr))
			})

		})
	})

	Describe("GetOutstandingAmount", func() {
		Describe("Positive case", func() {
			It("should return correct total unpaid payment without cache", func() {
				customerID := uuid.New()
				totalUnpaid := 5000000.0

				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				repo.EXPECT().GetCustomerByID(ctx, customerID).Return(&domain.Customer{}, nil)
				repo.EXPECT().GetTotalUnpaidPaymentOnActiveLoan(ctx, customerID).Return(totalUnpaid, nil)
				cache.EXPECT().Set(ctx, gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().LastActiveLoan(ctx, customerID).Return(&domain.Loan{
					LoanID: customerID,
				}, nil)

				response, err := svc.GetOutstandingBalance(ctx, customerID)
				Expect(err).To(BeNil())
				Expect(response.OutstandingBalance).To(Equal(totalUnpaid))
			})

			It("should return correct total unpaid payment with cache", func() {
				customerID := uuid.New()
				totalUnpaid := 5000000.0
				cacheRes := "{\"outstanding_balance\":5000000}"

				cache.EXPECT().Get(ctx, gomock.Any()).Return(cacheRes, nil)

				response, err := svc.GetOutstandingBalance(ctx, customerID)
				Expect(err).To(BeNil())
				Expect(response.OutstandingBalance).To(Equal(totalUnpaid))
			})

		})

		Describe("Negative case", func() {
			It("when error getting cache", func() {
				customerID := uuid.New()
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, someErr)
				_, err := svc.GetOutstandingBalance(ctx, customerID)
				Expect(err).To(Equal(someErr))
			})

			It("when error getting customer by id", func() {
				customerID := uuid.New()
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				repo.EXPECT().GetCustomerByID(ctx, customerID).Return(nil, someErr)
				_, err := svc.GetOutstandingBalance(ctx, customerID)
				Expect(err).To(Equal(someErr))
			})

			It("when error getting total unpaid payment on active loan", func() {
				customerID := uuid.New()
				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				repo.EXPECT().GetCustomerByID(ctx, customerID).Return(&domain.Customer{}, nil)
				repo.EXPECT().GetTotalUnpaidPaymentOnActiveLoan(ctx, customerID).Return(float64(0), someErr)
				repo.EXPECT().LastActiveLoan(ctx, customerID).Return(&domain.Loan{
					LoanID: customerID,
				}, nil)

				_, err := svc.GetOutstandingBalance(ctx, customerID)
				Expect(err).To(Equal(someErr))
			})

			It("when error getting last active loan", func() {
				customerID := uuid.New()

				cache.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)
				repo.EXPECT().GetCustomerByID(ctx, customerID).Return(&domain.Customer{}, nil)
				repo.EXPECT().LastActiveLoan(ctx, customerID).Return(nil, someErr)

				_, err := svc.GetOutstandingBalance(ctx, customerID)
				Expect(err).To(Equal(someErr))
			})
		})
	})
})
