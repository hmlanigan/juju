// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/containerbroker (interfaces: State)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	provisioner "github.com/juju/juju/api/agent/provisioner"
	network "github.com/juju/juju/core/network"
	network0 "github.com/juju/juju/network"
	params "github.com/juju/juju/rpc/params"
	names "github.com/juju/names/v4"
)

// MockState is a mock of State interface.
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState.
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance.
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// ContainerConfig mocks base method.
func (m *MockState) ContainerConfig() (params.ContainerConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerConfig")
	ret0, _ := ret[0].(params.ContainerConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerConfig indicates an expected call of ContainerConfig.
func (mr *MockStateMockRecorder) ContainerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerConfig", reflect.TypeOf((*MockState)(nil).ContainerConfig))
}

// ContainerManagerConfig mocks base method.
func (m *MockState) ContainerManagerConfig(arg0 params.ContainerManagerConfigParams) (params.ContainerManagerConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerManagerConfig", arg0)
	ret0, _ := ret[0].(params.ContainerManagerConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerManagerConfig indicates an expected call of ContainerManagerConfig.
func (mr *MockStateMockRecorder) ContainerManagerConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerManagerConfig", reflect.TypeOf((*MockState)(nil).ContainerManagerConfig), arg0)
}

// GetContainerProfileInfo mocks base method.
func (m *MockState) GetContainerProfileInfo(arg0 names.MachineTag) ([]*provisioner.LXDProfileResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContainerProfileInfo", arg0)
	ret0, _ := ret[0].([]*provisioner.LXDProfileResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContainerProfileInfo indicates an expected call of GetContainerProfileInfo.
func (mr *MockStateMockRecorder) GetContainerProfileInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContainerProfileInfo", reflect.TypeOf((*MockState)(nil).GetContainerProfileInfo), arg0)
}

// HostChangesForContainer mocks base method.
func (m *MockState) HostChangesForContainer(arg0 names.MachineTag) ([]network0.DeviceToBridge, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HostChangesForContainer", arg0)
	ret0, _ := ret[0].([]network0.DeviceToBridge)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// HostChangesForContainer indicates an expected call of HostChangesForContainer.
func (mr *MockStateMockRecorder) HostChangesForContainer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HostChangesForContainer", reflect.TypeOf((*MockState)(nil).HostChangesForContainer), arg0)
}

// Machines mocks base method.
func (m *MockState) Machines(arg0 ...names.MachineTag) ([]provisioner.MachineResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Machines", varargs...)
	ret0, _ := ret[0].([]provisioner.MachineResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Machines indicates an expected call of Machines.
func (mr *MockStateMockRecorder) Machines(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machines", reflect.TypeOf((*MockState)(nil).Machines), arg0...)
}

// PrepareContainerInterfaceInfo mocks base method.
func (m *MockState) PrepareContainerInterfaceInfo(arg0 names.MachineTag) (network.InterfaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareContainerInterfaceInfo", arg0)
	ret0, _ := ret[0].(network.InterfaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareContainerInterfaceInfo indicates an expected call of PrepareContainerInterfaceInfo.
func (mr *MockStateMockRecorder) PrepareContainerInterfaceInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareContainerInterfaceInfo", reflect.TypeOf((*MockState)(nil).PrepareContainerInterfaceInfo), arg0)
}

// ReleaseContainerAddresses mocks base method.
func (m *MockState) ReleaseContainerAddresses(arg0 names.MachineTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseContainerAddresses", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseContainerAddresses indicates an expected call of ReleaseContainerAddresses.
func (mr *MockStateMockRecorder) ReleaseContainerAddresses(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseContainerAddresses", reflect.TypeOf((*MockState)(nil).ReleaseContainerAddresses), arg0)
}

// SetHostMachineNetworkConfig mocks base method.
func (m *MockState) SetHostMachineNetworkConfig(arg0 names.MachineTag, arg1 []params.NetworkConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHostMachineNetworkConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHostMachineNetworkConfig indicates an expected call of SetHostMachineNetworkConfig.
func (mr *MockStateMockRecorder) SetHostMachineNetworkConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHostMachineNetworkConfig", reflect.TypeOf((*MockState)(nil).SetHostMachineNetworkConfig), arg0, arg1)
}
