// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/worker/v4/dependency (interfaces: Getter)
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/dependency.go github.com/juju/worker/v4/dependency Getter
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockGetter is a mock of Getter interface.
type MockGetter struct {
	ctrl     *gomock.Controller
	recorder *MockGetterMockRecorder
}

// MockGetterMockRecorder is the mock recorder for MockGetter.
type MockGetterMockRecorder struct {
	mock *MockGetter
}

// NewMockGetter creates a new mock instance.
func NewMockGetter(ctrl *gomock.Controller) *MockGetter {
	mock := &MockGetter{ctrl: ctrl}
	mock.recorder = &MockGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetter) EXPECT() *MockGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockGetter) Get(arg0 string, arg1 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockGetterMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGetter)(nil).Get), arg0, arg1)
}
