// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	entity "password_keeper/internal/common/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthorizationRepository is a mock of AuthorizationRepository interface.
type MockAuthorizationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationRepositoryMockRecorder
}

// MockAuthorizationRepositoryMockRecorder is the mock recorder for MockAuthorizationRepository.
type MockAuthorizationRepositoryMockRecorder struct {
	mock *MockAuthorizationRepository
}

// NewMockAuthorizationRepository creates a new mock instance.
func NewMockAuthorizationRepository(ctrl *gomock.Controller) *MockAuthorizationRepository {
	mock := &MockAuthorizationRepository{ctrl: ctrl}
	mock.recorder = &MockAuthorizationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizationRepository) EXPECT() *MockAuthorizationRepositoryMockRecorder {
	return m.recorder
}

// GetUserFromDB mocks base method.
func (m *MockAuthorizationRepository) GetUserFromDB(ctx context.Context, user entity.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFromDB", ctx, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFromDB indicates an expected call of GetUserFromDB.
func (mr *MockAuthorizationRepositoryMockRecorder) GetUserFromDB(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFromDB", reflect.TypeOf((*MockAuthorizationRepository)(nil).GetUserFromDB), ctx, user)
}

// SetUserDB mocks base method.
func (m *MockAuthorizationRepository) SetUserDB(ctx context.Context, user entity.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserDB", ctx, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetUserDB indicates an expected call of SetUserDB.
func (mr *MockAuthorizationRepositoryMockRecorder) SetUserDB(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserDB", reflect.TypeOf((*MockAuthorizationRepository)(nil).SetUserDB), ctx, user)
}

// Validate mocks base method.
func (m *MockAuthorizationRepository) Validate(ctx context.Context, username string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", ctx, username)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockAuthorizationRepositoryMockRecorder) Validate(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockAuthorizationRepository)(nil).Validate), ctx, username)
}

// MockDataRepositoryInterface is a mock of DataRepositoryInterface interface.
type MockDataRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDataRepositoryInterfaceMockRecorder
}

// MockDataRepositoryInterfaceMockRecorder is the mock recorder for MockDataRepositoryInterface.
type MockDataRepositoryInterfaceMockRecorder struct {
	mock *MockDataRepositoryInterface
}

// NewMockDataRepositoryInterface creates a new mock instance.
func NewMockDataRepositoryInterface(ctrl *gomock.Controller) *MockDataRepositoryInterface {
	mock := &MockDataRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockDataRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataRepositoryInterface) EXPECT() *MockDataRepositoryInterfaceMockRecorder {
	return m.recorder
}

// DeleteRepData mocks base method.
func (m *MockDataRepositoryInterface) DeleteRepData(id int, eventType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRepData", id, eventType)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRepData indicates an expected call of DeleteRepData.
func (mr *MockDataRepositoryInterfaceMockRecorder) DeleteRepData(id, eventType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRepData", reflect.TypeOf((*MockDataRepositoryInterface)(nil).DeleteRepData), id, eventType)
}

// GetRepData mocks base method.
func (m *MockDataRepositoryInterface) GetRepData(id int, eventType string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepData", id, eventType)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepData indicates an expected call of GetRepData.
func (mr *MockDataRepositoryInterfaceMockRecorder) GetRepData(id, eventType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepData", reflect.TypeOf((*MockDataRepositoryInterface)(nil).GetRepData), id, eventType)
}

// SetRepData mocks base method.
func (m *MockDataRepositoryInterface) SetRepData(id int, text []byte, eventType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRepData", id, text, eventType)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRepData indicates an expected call of SetRepData.
func (mr *MockDataRepositoryInterfaceMockRecorder) SetRepData(id, text, eventType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRepData", reflect.TypeOf((*MockDataRepositoryInterface)(nil).SetRepData), id, text, eventType)
}
