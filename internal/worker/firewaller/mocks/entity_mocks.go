// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/worker/firewaller (interfaces: Machine,Unit,Application)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	instance "github.com/juju/juju/core/instance"
	life "github.com/juju/juju/core/life"
	network "github.com/juju/juju/core/network"
	watcher "github.com/juju/juju/core/watcher"
	firewaller "github.com/juju/juju/internal/worker/firewaller"
	params "github.com/juju/juju/rpc/params"
	names "github.com/juju/names/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockMachine is a mock of Machine interface.
type MockMachine struct {
	ctrl     *gomock.Controller
	recorder *MockMachineMockRecorder
}

// MockMachineMockRecorder is the mock recorder for MockMachine.
type MockMachineMockRecorder struct {
	mock *MockMachine
}

// NewMockMachine creates a new mock instance.
func NewMockMachine(ctrl *gomock.Controller) *MockMachine {
	mock := &MockMachine{ctrl: ctrl}
	mock.recorder = &MockMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachine) EXPECT() *MockMachineMockRecorder {
	return m.recorder
}

// InstanceId mocks base method.
func (m *MockMachine) InstanceId() (instance.Id, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceId")
	ret0, _ := ret[0].(instance.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceId indicates an expected call of InstanceId.
func (mr *MockMachineMockRecorder) InstanceId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceId", reflect.TypeOf((*MockMachine)(nil).InstanceId))
}

// IsManual mocks base method.
func (m *MockMachine) IsManual() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsManual")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsManual indicates an expected call of IsManual.
func (mr *MockMachineMockRecorder) IsManual() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsManual", reflect.TypeOf((*MockMachine)(nil).IsManual))
}

// Life mocks base method.
func (m *MockMachine) Life() life.Value {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Life")
	ret0, _ := ret[0].(life.Value)
	return ret0
}

// Life indicates an expected call of Life.
func (mr *MockMachineMockRecorder) Life() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Life", reflect.TypeOf((*MockMachine)(nil).Life))
}

// OpenedMachinePortRanges mocks base method.
func (m *MockMachine) OpenedMachinePortRanges() (map[names.UnitTag]network.GroupedPortRanges, map[names.UnitTag]network.GroupedPortRanges, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenedMachinePortRanges")
	ret0, _ := ret[0].(map[names.UnitTag]network.GroupedPortRanges)
	ret1, _ := ret[1].(map[names.UnitTag]network.GroupedPortRanges)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// OpenedMachinePortRanges indicates an expected call of OpenedMachinePortRanges.
func (mr *MockMachineMockRecorder) OpenedMachinePortRanges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenedMachinePortRanges", reflect.TypeOf((*MockMachine)(nil).OpenedMachinePortRanges))
}

// Tag mocks base method.
func (m *MockMachine) Tag() names.MachineTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names.MachineTag)
	return ret0
}

// Tag indicates an expected call of Tag.
func (mr *MockMachineMockRecorder) Tag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockMachine)(nil).Tag))
}

// WatchUnits mocks base method.
func (m *MockMachine) WatchUnits() (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchUnits")
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchUnits indicates an expected call of WatchUnits.
func (mr *MockMachineMockRecorder) WatchUnits() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchUnits", reflect.TypeOf((*MockMachine)(nil).WatchUnits))
}

// MockUnit is a mock of Unit interface.
type MockUnit struct {
	ctrl     *gomock.Controller
	recorder *MockUnitMockRecorder
}

// MockUnitMockRecorder is the mock recorder for MockUnit.
type MockUnitMockRecorder struct {
	mock *MockUnit
}

// NewMockUnit creates a new mock instance.
func NewMockUnit(ctrl *gomock.Controller) *MockUnit {
	mock := &MockUnit{ctrl: ctrl}
	mock.recorder = &MockUnitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnit) EXPECT() *MockUnitMockRecorder {
	return m.recorder
}

// Application mocks base method.
func (m *MockUnit) Application() (firewaller.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application")
	ret0, _ := ret[0].(firewaller.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Application indicates an expected call of Application.
func (mr *MockUnitMockRecorder) Application() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockUnit)(nil).Application))
}

// AssignedMachine mocks base method.
func (m *MockUnit) AssignedMachine() (names.MachineTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignedMachine")
	ret0, _ := ret[0].(names.MachineTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignedMachine indicates an expected call of AssignedMachine.
func (mr *MockUnitMockRecorder) AssignedMachine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignedMachine", reflect.TypeOf((*MockUnit)(nil).AssignedMachine))
}

// Life mocks base method.
func (m *MockUnit) Life() life.Value {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Life")
	ret0, _ := ret[0].(life.Value)
	return ret0
}

// Life indicates an expected call of Life.
func (mr *MockUnitMockRecorder) Life() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Life", reflect.TypeOf((*MockUnit)(nil).Life))
}

// Name mocks base method.
func (m *MockUnit) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockUnitMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockUnit)(nil).Name))
}

// Refresh mocks base method.
func (m *MockUnit) Refresh() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh.
func (mr *MockUnitMockRecorder) Refresh() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockUnit)(nil).Refresh))
}

// Tag mocks base method.
func (m *MockUnit) Tag() names.UnitTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names.UnitTag)
	return ret0
}

// Tag indicates an expected call of Tag.
func (mr *MockUnitMockRecorder) Tag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockUnit)(nil).Tag))
}

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// ExposeInfo mocks base method.
func (m *MockApplication) ExposeInfo() (bool, map[string]params.ExposedEndpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExposeInfo")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(map[string]params.ExposedEndpoint)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ExposeInfo indicates an expected call of ExposeInfo.
func (mr *MockApplicationMockRecorder) ExposeInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExposeInfo", reflect.TypeOf((*MockApplication)(nil).ExposeInfo))
}

// Name mocks base method.
func (m *MockApplication) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockApplicationMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockApplication)(nil).Name))
}

// Tag mocks base method.
func (m *MockApplication) Tag() names.ApplicationTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names.ApplicationTag)
	return ret0
}

// Tag indicates an expected call of Tag.
func (mr *MockApplicationMockRecorder) Tag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockApplication)(nil).Tag))
}

// Watch mocks base method.
func (m *MockApplication) Watch() (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch")
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch.
func (mr *MockApplicationMockRecorder) Watch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockApplication)(nil).Watch))
}
