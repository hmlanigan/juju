// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/common (interfaces: UpgradeSeriesBackend,UpgradeSeriesMachine,UpgradeSeriesUnit)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	common "github.com/juju/juju/apiserver/common"
	model "github.com/juju/juju/core/model"
	state "github.com/juju/juju/state"
	names_v2 "gopkg.in/juju/names.v2"
	reflect "reflect"
)

// MockUpgradeSeriesBackend is a mock of UpgradeSeriesBackend interface
type MockUpgradeSeriesBackend struct {
	ctrl     *gomock.Controller
	recorder *MockUpgradeSeriesBackendMockRecorder
}

// MockUpgradeSeriesBackendMockRecorder is the mock recorder for MockUpgradeSeriesBackend
type MockUpgradeSeriesBackendMockRecorder struct {
	mock *MockUpgradeSeriesBackend
}

// NewMockUpgradeSeriesBackend creates a new mock instance
func NewMockUpgradeSeriesBackend(ctrl *gomock.Controller) *MockUpgradeSeriesBackend {
	mock := &MockUpgradeSeriesBackend{ctrl: ctrl}
	mock.recorder = &MockUpgradeSeriesBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUpgradeSeriesBackend) EXPECT() *MockUpgradeSeriesBackendMockRecorder {
	return m.recorder
}

// Machine mocks base method
func (m *MockUpgradeSeriesBackend) Machine(arg0 string) (common.UpgradeSeriesMachine, error) {
	ret := m.ctrl.Call(m, "Machine", arg0)
	ret0, _ := ret[0].(common.UpgradeSeriesMachine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Machine indicates an expected call of Machine
func (mr *MockUpgradeSeriesBackendMockRecorder) Machine(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machine", reflect.TypeOf((*MockUpgradeSeriesBackend)(nil).Machine), arg0)
}

// Unit mocks base method
func (m *MockUpgradeSeriesBackend) Unit(arg0 string) (common.UpgradeSeriesUnit, error) {
	ret := m.ctrl.Call(m, "Unit", arg0)
	ret0, _ := ret[0].(common.UpgradeSeriesUnit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unit indicates an expected call of Unit
func (mr *MockUpgradeSeriesBackendMockRecorder) Unit(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unit", reflect.TypeOf((*MockUpgradeSeriesBackend)(nil).Unit), arg0)
}

// MockUpgradeSeriesMachine is a mock of UpgradeSeriesMachine interface
type MockUpgradeSeriesMachine struct {
	ctrl     *gomock.Controller
	recorder *MockUpgradeSeriesMachineMockRecorder
}

// MockUpgradeSeriesMachineMockRecorder is the mock recorder for MockUpgradeSeriesMachine
type MockUpgradeSeriesMachineMockRecorder struct {
	mock *MockUpgradeSeriesMachine
}

// NewMockUpgradeSeriesMachine creates a new mock instance
func NewMockUpgradeSeriesMachine(ctrl *gomock.Controller) *MockUpgradeSeriesMachine {
	mock := &MockUpgradeSeriesMachine{ctrl: ctrl}
	mock.recorder = &MockUpgradeSeriesMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUpgradeSeriesMachine) EXPECT() *MockUpgradeSeriesMachineMockRecorder {
	return m.recorder
}

// MachineUpgradeSeriesStatus mocks base method
func (m *MockUpgradeSeriesMachine) MachineUpgradeSeriesStatus() (model.UpgradeSeriesStatus, error) {
	ret := m.ctrl.Call(m, "MachineUpgradeSeriesStatus")
	ret0, _ := ret[0].(model.UpgradeSeriesStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MachineUpgradeSeriesStatus indicates an expected call of MachineUpgradeSeriesStatus
func (mr *MockUpgradeSeriesMachineMockRecorder) MachineUpgradeSeriesStatus() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MachineUpgradeSeriesStatus", reflect.TypeOf((*MockUpgradeSeriesMachine)(nil).MachineUpgradeSeriesStatus))
}

// SetMachineUpgradeSeriesStatus mocks base method
func (m *MockUpgradeSeriesMachine) SetMachineUpgradeSeriesStatus(arg0 model.UpgradeSeriesStatus) error {
	ret := m.ctrl.Call(m, "SetMachineUpgradeSeriesStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMachineUpgradeSeriesStatus indicates an expected call of SetMachineUpgradeSeriesStatus
func (mr *MockUpgradeSeriesMachineMockRecorder) SetMachineUpgradeSeriesStatus(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMachineUpgradeSeriesStatus", reflect.TypeOf((*MockUpgradeSeriesMachine)(nil).SetMachineUpgradeSeriesStatus), arg0)
}

// StartUnitUpgradeSeriesCompletionPhase mocks base method
func (m *MockUpgradeSeriesMachine) StartUnitUpgradeSeriesCompletionPhase() error {
	ret := m.ctrl.Call(m, "StartUnitUpgradeSeriesCompletionPhase")
	ret0, _ := ret[0].(error)
	return ret0
}

// StartUnitUpgradeSeriesCompletionPhase indicates an expected call of StartUnitUpgradeSeriesCompletionPhase
func (mr *MockUpgradeSeriesMachineMockRecorder) StartUnitUpgradeSeriesCompletionPhase() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartUnitUpgradeSeriesCompletionPhase", reflect.TypeOf((*MockUpgradeSeriesMachine)(nil).StartUnitUpgradeSeriesCompletionPhase))
}

// Units mocks base method
func (m *MockUpgradeSeriesMachine) Units() ([]common.UpgradeSeriesUnit, error) {
	ret := m.ctrl.Call(m, "Units")
	ret0, _ := ret[0].([]common.UpgradeSeriesUnit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Units indicates an expected call of Units
func (mr *MockUpgradeSeriesMachineMockRecorder) Units() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Units", reflect.TypeOf((*MockUpgradeSeriesMachine)(nil).Units))
}

// WatchUpgradeSeriesNotifications mocks base method
func (m *MockUpgradeSeriesMachine) WatchUpgradeSeriesNotifications() (state.NotifyWatcher, error) {
	ret := m.ctrl.Call(m, "WatchUpgradeSeriesNotifications")
	ret0, _ := ret[0].(state.NotifyWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchUpgradeSeriesNotifications indicates an expected call of WatchUpgradeSeriesNotifications
func (mr *MockUpgradeSeriesMachineMockRecorder) WatchUpgradeSeriesNotifications() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchUpgradeSeriesNotifications", reflect.TypeOf((*MockUpgradeSeriesMachine)(nil).WatchUpgradeSeriesNotifications))
}

// MockUpgradeSeriesUnit is a mock of UpgradeSeriesUnit interface
type MockUpgradeSeriesUnit struct {
	ctrl     *gomock.Controller
	recorder *MockUpgradeSeriesUnitMockRecorder
}

// MockUpgradeSeriesUnitMockRecorder is the mock recorder for MockUpgradeSeriesUnit
type MockUpgradeSeriesUnitMockRecorder struct {
	mock *MockUpgradeSeriesUnit
}

// NewMockUpgradeSeriesUnit creates a new mock instance
func NewMockUpgradeSeriesUnit(ctrl *gomock.Controller) *MockUpgradeSeriesUnit {
	mock := &MockUpgradeSeriesUnit{ctrl: ctrl}
	mock.recorder = &MockUpgradeSeriesUnitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUpgradeSeriesUnit) EXPECT() *MockUpgradeSeriesUnitMockRecorder {
	return m.recorder
}

// AssignedMachineId mocks base method
func (m *MockUpgradeSeriesUnit) AssignedMachineId() (string, error) {
	ret := m.ctrl.Call(m, "AssignedMachineId")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignedMachineId indicates an expected call of AssignedMachineId
func (mr *MockUpgradeSeriesUnitMockRecorder) AssignedMachineId() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignedMachineId", reflect.TypeOf((*MockUpgradeSeriesUnit)(nil).AssignedMachineId))
}

// SetUpgradeSeriesStatus mocks base method
func (m *MockUpgradeSeriesUnit) SetUpgradeSeriesStatus(arg0 model.UpgradeSeriesStatus) error {
	ret := m.ctrl.Call(m, "SetUpgradeSeriesStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUpgradeSeriesStatus indicates an expected call of SetUpgradeSeriesStatus
func (mr *MockUpgradeSeriesUnitMockRecorder) SetUpgradeSeriesStatus(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUpgradeSeriesStatus", reflect.TypeOf((*MockUpgradeSeriesUnit)(nil).SetUpgradeSeriesStatus), arg0)
}

// Tag mocks base method
func (m *MockUpgradeSeriesUnit) Tag() names_v2.Tag {
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names_v2.Tag)
	return ret0
}

// Tag indicates an expected call of Tag
func (mr *MockUpgradeSeriesUnitMockRecorder) Tag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockUpgradeSeriesUnit)(nil).Tag))
}

// UpgradeSeriesStatus mocks base method
func (m *MockUpgradeSeriesUnit) UpgradeSeriesStatus() (model.UpgradeSeriesStatus, error) {
	ret := m.ctrl.Call(m, "UpgradeSeriesStatus")
	ret0, _ := ret[0].(model.UpgradeSeriesStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpgradeSeriesStatus indicates an expected call of UpgradeSeriesStatus
func (mr *MockUpgradeSeriesUnitMockRecorder) UpgradeSeriesStatus() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeSeriesStatus", reflect.TypeOf((*MockUpgradeSeriesUnit)(nil).UpgradeSeriesStatus))
}
