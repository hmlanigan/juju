// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/application (interfaces: StateModel)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/juju/juju/environs/config"
)

// MockStateModel is a mock of StateModel interface.
type MockStateModel struct {
	ctrl     *gomock.Controller
	recorder *MockStateModelMockRecorder
}

// MockStateModelMockRecorder is the mock recorder for MockStateModel.
type MockStateModelMockRecorder struct {
	mock *MockStateModel
}

// NewMockStateModel creates a new mock instance.
func NewMockStateModel(ctrl *gomock.Controller) *MockStateModel {
	mock := &MockStateModel{ctrl: ctrl}
	mock.recorder = &MockStateModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateModel) EXPECT() *MockStateModelMockRecorder {
	return m.recorder
}

// ModelConfig mocks base method.
func (m *MockStateModel) ModelConfig() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig.
func (mr *MockStateModelMockRecorder) ModelConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockStateModel)(nil).ModelConfig))
}
