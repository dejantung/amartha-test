// Code generated by MockGen. DO NOT EDIT.
// Source: billing-engine/internal/billing/repository (interfaces: BillingRepositoryProvider)
//
// Generated by this command:
//
//	mockgen -destination=../mocks/mock_billing_repository.go -package=mocks billing-engine/internal/billing/repository BillingRepositoryProvider
//

// Package mocks is a generated GoMock package.
package mocks

import (
	domain "billing-engine/internal/billing/domain"
	context "context"
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockBillingRepositoryProvider is a mock of BillingRepositoryProvider interface.
type MockBillingRepositoryProvider struct {
	ctrl     *gomock.Controller
	recorder *MockBillingRepositoryProviderMockRecorder
}

// MockBillingRepositoryProviderMockRecorder is the mock recorder for MockBillingRepositoryProvider.
type MockBillingRepositoryProviderMockRecorder struct {
	mock *MockBillingRepositoryProvider
}

// NewMockBillingRepositoryProvider creates a new mock instance.
func NewMockBillingRepositoryProvider(ctrl *gomock.Controller) *MockBillingRepositoryProvider {
	mock := &MockBillingRepositoryProvider{ctrl: ctrl}
	mock.recorder = &MockBillingRepositoryProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBillingRepositoryProvider) EXPECT() *MockBillingRepositoryProviderMockRecorder {
	return m.recorder
}

// CreateCustomer mocks base method.
func (m *MockBillingRepositoryProvider) CreateCustomer(arg0 context.Context, arg1 []domain.Customer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCustomer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCustomer indicates an expected call of CreateCustomer.
func (mr *MockBillingRepositoryProviderMockRecorder) CreateCustomer(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).CreateCustomer), arg0, arg1)
}

// CreateLoan mocks base method.
func (m *MockBillingRepositoryProvider) CreateLoan(arg0 context.Context, arg1 domain.Loan) (*domain.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLoan", arg0, arg1)
	ret0, _ := ret[0].(*domain.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLoan indicates an expected call of CreateLoan.
func (mr *MockBillingRepositoryProviderMockRecorder) CreateLoan(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoan", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).CreateLoan), arg0, arg1)
}

// GetCustomer mocks base method.
func (m *MockBillingRepositoryProvider) GetCustomer(arg0 context.Context) ([]domain.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomer", arg0)
	ret0, _ := ret[0].([]domain.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCustomer indicates an expected call of GetCustomer.
func (mr *MockBillingRepositoryProviderMockRecorder) GetCustomer(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomer", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).GetCustomer), arg0)
}

// GetCustomerByID mocks base method.
func (m *MockBillingRepositoryProvider) GetCustomerByID(arg0 context.Context, arg1 uuid.UUID) (*domain.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomerByID", arg0, arg1)
	ret0, _ := ret[0].(*domain.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCustomerByID indicates an expected call of GetCustomerByID.
func (mr *MockBillingRepositoryProviderMockRecorder) GetCustomerByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomerByID", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).GetCustomerByID), arg0, arg1)
}

// GetLoanByIDAndCustomerID mocks base method.
func (m *MockBillingRepositoryProvider) GetLoanByIDAndCustomerID(arg0 context.Context, arg1, arg2 uuid.UUID) (*domain.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoanByIDAndCustomerID", arg0, arg1, arg2)
	ret0, _ := ret[0].(*domain.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoanByIDAndCustomerID indicates an expected call of GetLoanByIDAndCustomerID.
func (mr *MockBillingRepositoryProviderMockRecorder) GetLoanByIDAndCustomerID(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoanByIDAndCustomerID", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).GetLoanByIDAndCustomerID), arg0, arg1, arg2)
}

// GetLoanByScheduleID mocks base method.
func (m *MockBillingRepositoryProvider) GetLoanByScheduleID(arg0 context.Context, arg1 uuid.UUID) (*domain.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoanByScheduleID", arg0, arg1)
	ret0, _ := ret[0].(*domain.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoanByScheduleID indicates an expected call of GetLoanByScheduleID.
func (mr *MockBillingRepositoryProviderMockRecorder) GetLoanByScheduleID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoanByScheduleID", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).GetLoanByScheduleID), arg0, arg1)
}

// GetSchedule mocks base method.
func (m *MockBillingRepositoryProvider) GetSchedule(arg0 context.Context, arg1, arg2 uuid.UUID) ([]domain.Schedule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchedule", arg0, arg1, arg2)
	ret0, _ := ret[0].([]domain.Schedule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchedule indicates an expected call of GetSchedule.
func (mr *MockBillingRepositoryProviderMockRecorder) GetSchedule(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchedule", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).GetSchedule), arg0, arg1, arg2)
}

// GetTotalUnpaidPaymentOnActiveLoan mocks base method.
func (m *MockBillingRepositoryProvider) GetTotalUnpaidPaymentOnActiveLoan(arg0 context.Context, arg1 uuid.UUID) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalUnpaidPaymentOnActiveLoan", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalUnpaidPaymentOnActiveLoan indicates an expected call of GetTotalUnpaidPaymentOnActiveLoan.
func (mr *MockBillingRepositoryProviderMockRecorder) GetTotalUnpaidPaymentOnActiveLoan(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalUnpaidPaymentOnActiveLoan", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).GetTotalUnpaidPaymentOnActiveLoan), arg0, arg1)
}

// GetUnpaidAndMissPaymentUntil mocks base method.
func (m *MockBillingRepositoryProvider) GetUnpaidAndMissPaymentUntil(arg0 context.Context, arg1 uuid.UUID, arg2 time.Time) ([]domain.Schedule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnpaidAndMissPaymentUntil", arg0, arg1, arg2)
	ret0, _ := ret[0].([]domain.Schedule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnpaidAndMissPaymentUntil indicates an expected call of GetUnpaidAndMissPaymentUntil.
func (mr *MockBillingRepositoryProviderMockRecorder) GetUnpaidAndMissPaymentUntil(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnpaidAndMissPaymentUntil", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).GetUnpaidAndMissPaymentUntil), arg0, arg1, arg2)
}

// LastActiveLoan mocks base method.
func (m *MockBillingRepositoryProvider) LastActiveLoan(arg0 context.Context, arg1 uuid.UUID) (*domain.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastActiveLoan", arg0, arg1)
	ret0, _ := ret[0].(*domain.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastActiveLoan indicates an expected call of LastActiveLoan.
func (mr *MockBillingRepositoryProviderMockRecorder) LastActiveLoan(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastActiveLoan", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).LastActiveLoan), arg0, arg1)
}

// UpdateSchedulePayment mocks base method.
func (m *MockBillingRepositoryProvider) UpdateSchedulePayment(arg0 context.Context, arg1 *domain.Schedule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSchedulePayment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSchedulePayment indicates an expected call of UpdateSchedulePayment.
func (mr *MockBillingRepositoryProviderMockRecorder) UpdateSchedulePayment(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSchedulePayment", reflect.TypeOf((*MockBillingRepositoryProvider)(nil).UpdateSchedulePayment), arg0, arg1)
}
