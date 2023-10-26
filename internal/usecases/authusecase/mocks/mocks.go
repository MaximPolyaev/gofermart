// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entities "github.com/MaximPolyaev/gofermart/internal/entities"
	gomock "github.com/golang/mock/gomock"
)

// Mockstorage is a mock of storage interface.
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
}

// MockstorageMockRecorder is the mock recorder for Mockstorage.
type MockstorageMockRecorder struct {
	mock *Mockstorage
}

// NewMockstorage creates a new mock instance.
func NewMockstorage(ctrl *gomock.Controller) *Mockstorage {
	mock := &Mockstorage{ctrl: ctrl}
	mock.recorder = &MockstorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockstorage) EXPECT() *MockstorageMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *Mockstorage) CreateUser(ctx context.Context, user *entities.UserWithPassword) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockstorageMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*Mockstorage)(nil).CreateUser), ctx, user)
}

// FindUserWithPassword mocks base method.
func (m *Mockstorage) FindUserWithPassword(ctx context.Context, login string) (*entities.UserWithPassword, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserWithPassword", ctx, login)
	ret0, _ := ret[0].(*entities.UserWithPassword)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserWithPassword indicates an expected call of FindUserWithPassword.
func (mr *MockstorageMockRecorder) FindUserWithPassword(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserWithPassword", reflect.TypeOf((*Mockstorage)(nil).FindUserWithPassword), ctx, login)
}
