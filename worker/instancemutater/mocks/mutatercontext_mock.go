// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/instancemutater (interfaces: MutaterContext)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	instancemutater "github.com/juju/juju/api/agent/instancemutater"
	environs "github.com/juju/juju/environs"
	instancemutater0 "github.com/juju/juju/worker/instancemutater"
	names "github.com/juju/names/v4"
	worker "github.com/juju/worker/v3"
)

// MockMutaterContext is a mock of MutaterContext interface.
type MockMutaterContext struct {
	ctrl     *gomock.Controller
	recorder *MockMutaterContextMockRecorder
}

// MockMutaterContextMockRecorder is the mock recorder for MockMutaterContext.
type MockMutaterContextMockRecorder struct {
	mock *MockMutaterContext
}

// NewMockMutaterContext creates a new mock instance.
func NewMockMutaterContext(ctrl *gomock.Controller) *MockMutaterContext {
	mock := &MockMutaterContext{ctrl: ctrl}
	mock.recorder = &MockMutaterContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMutaterContext) EXPECT() *MockMutaterContextMockRecorder {
	return m.recorder
}

// KillWithError mocks base method.
func (m *MockMutaterContext) KillWithError(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "KillWithError", arg0)
}

// KillWithError indicates an expected call of KillWithError.
func (mr *MockMutaterContextMockRecorder) KillWithError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KillWithError", reflect.TypeOf((*MockMutaterContext)(nil).KillWithError), arg0)
}

// add mocks base method.
func (m *MockMutaterContext) add(arg0 worker.Worker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// add indicates an expected call of add.
func (mr *MockMutaterContextMockRecorder) add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "add", reflect.TypeOf((*MockMutaterContext)(nil).add), arg0)
}

// dying mocks base method.
func (m *MockMutaterContext) dying() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "dying")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// dying indicates an expected call of dying.
func (mr *MockMutaterContextMockRecorder) dying() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "dying", reflect.TypeOf((*MockMutaterContext)(nil).dying))
}

// errDying mocks base method.
func (m *MockMutaterContext) errDying() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "errDying")
	ret0, _ := ret[0].(error)
	return ret0
}

// errDying indicates an expected call of errDying.
func (mr *MockMutaterContextMockRecorder) errDying() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "errDying", reflect.TypeOf((*MockMutaterContext)(nil).errDying))
}

// getBroker mocks base method.
func (m *MockMutaterContext) getBroker() environs.LXDProfiler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getBroker")
	ret0, _ := ret[0].(environs.LXDProfiler)
	return ret0
}

// getBroker indicates an expected call of getBroker.
func (mr *MockMutaterContextMockRecorder) getBroker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getBroker", reflect.TypeOf((*MockMutaterContext)(nil).getBroker))
}

// getMachine mocks base method.
func (m *MockMutaterContext) getMachine(arg0 names.MachineTag) (instancemutater.MutaterMachine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getMachine", arg0)
	ret0, _ := ret[0].(instancemutater.MutaterMachine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getMachine indicates an expected call of getMachine.
func (mr *MockMutaterContextMockRecorder) getMachine(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getMachine", reflect.TypeOf((*MockMutaterContext)(nil).getMachine), arg0)
}

// getRequiredLXDProfiles mocks base method.
func (m *MockMutaterContext) getRequiredLXDProfiles(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getRequiredLXDProfiles", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// getRequiredLXDProfiles indicates an expected call of getRequiredLXDProfiles.
func (mr *MockMutaterContextMockRecorder) getRequiredLXDProfiles(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getRequiredLXDProfiles", reflect.TypeOf((*MockMutaterContext)(nil).getRequiredLXDProfiles), arg0)
}

// newMachineContext mocks base method.
func (m *MockMutaterContext) newMachineContext() instancemutater0.MachineContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "newMachineContext")
	ret0, _ := ret[0].(instancemutater0.MachineContext)
	return ret0
}

// newMachineContext indicates an expected call of newMachineContext.
func (mr *MockMutaterContextMockRecorder) newMachineContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "newMachineContext", reflect.TypeOf((*MockMutaterContext)(nil).newMachineContext))
}
