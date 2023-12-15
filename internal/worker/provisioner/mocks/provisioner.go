// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/worker/provisioner (interfaces: ContainerMachine,ContainerMachineGetter,ContainerProvisionerAPI,ControllerAPI,MachinesAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	provisioner "github.com/juju/juju/api/agent/provisioner"
	controller "github.com/juju/juju/controller"
	instance "github.com/juju/juju/core/instance"
	life "github.com/juju/juju/core/life"
	network "github.com/juju/juju/core/network"
	watcher "github.com/juju/juju/core/watcher"
	config "github.com/juju/juju/environs/config"
	network0 "github.com/juju/juju/internal/network"
	provisioner0 "github.com/juju/juju/internal/worker/provisioner"
	params "github.com/juju/juju/rpc/params"
	names "github.com/juju/names/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockContainerMachine is a mock of ContainerMachine interface.
type MockContainerMachine struct {
	ctrl     *gomock.Controller
	recorder *MockContainerMachineMockRecorder
}

// MockContainerMachineMockRecorder is the mock recorder for MockContainerMachine.
type MockContainerMachineMockRecorder struct {
	mock *MockContainerMachine
}

// NewMockContainerMachine creates a new mock instance.
func NewMockContainerMachine(ctrl *gomock.Controller) *MockContainerMachine {
	mock := &MockContainerMachine{ctrl: ctrl}
	mock.recorder = &MockContainerMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContainerMachine) EXPECT() *MockContainerMachineMockRecorder {
	return m.recorder
}

// AvailabilityZone mocks base method.
func (m *MockContainerMachine) AvailabilityZone() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvailabilityZone")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AvailabilityZone indicates an expected call of AvailabilityZone.
func (mr *MockContainerMachineMockRecorder) AvailabilityZone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailabilityZone", reflect.TypeOf((*MockContainerMachine)(nil).AvailabilityZone))
}

// Life mocks base method.
func (m *MockContainerMachine) Life() life.Value {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Life")
	ret0, _ := ret[0].(life.Value)
	return ret0
}

// Life indicates an expected call of Life.
func (mr *MockContainerMachineMockRecorder) Life() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Life", reflect.TypeOf((*MockContainerMachine)(nil).Life))
}

// SupportedContainers mocks base method.
func (m *MockContainerMachine) SupportedContainers() ([]instance.ContainerType, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportedContainers")
	ret0, _ := ret[0].([]instance.ContainerType)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SupportedContainers indicates an expected call of SupportedContainers.
func (mr *MockContainerMachineMockRecorder) SupportedContainers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportedContainers", reflect.TypeOf((*MockContainerMachine)(nil).SupportedContainers))
}

// WatchContainers mocks base method.
func (m *MockContainerMachine) WatchContainers(arg0 instance.ContainerType) (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchContainers", arg0)
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchContainers indicates an expected call of WatchContainers.
func (mr *MockContainerMachineMockRecorder) WatchContainers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchContainers", reflect.TypeOf((*MockContainerMachine)(nil).WatchContainers), arg0)
}

// MockContainerMachineGetter is a mock of ContainerMachineGetter interface.
type MockContainerMachineGetter struct {
	ctrl     *gomock.Controller
	recorder *MockContainerMachineGetterMockRecorder
}

// MockContainerMachineGetterMockRecorder is the mock recorder for MockContainerMachineGetter.
type MockContainerMachineGetterMockRecorder struct {
	mock *MockContainerMachineGetter
}

// NewMockContainerMachineGetter creates a new mock instance.
func NewMockContainerMachineGetter(ctrl *gomock.Controller) *MockContainerMachineGetter {
	mock := &MockContainerMachineGetter{ctrl: ctrl}
	mock.recorder = &MockContainerMachineGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContainerMachineGetter) EXPECT() *MockContainerMachineGetterMockRecorder {
	return m.recorder
}

// Machines mocks base method.
func (m *MockContainerMachineGetter) Machines(arg0 ...names.MachineTag) ([]provisioner0.ContainerMachineResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Machines", varargs...)
	ret0, _ := ret[0].([]provisioner0.ContainerMachineResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Machines indicates an expected call of Machines.
func (mr *MockContainerMachineGetterMockRecorder) Machines(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machines", reflect.TypeOf((*MockContainerMachineGetter)(nil).Machines), arg0...)
}

// MockContainerProvisionerAPI is a mock of ContainerProvisionerAPI interface.
type MockContainerProvisionerAPI struct {
	ctrl     *gomock.Controller
	recorder *MockContainerProvisionerAPIMockRecorder
}

// MockContainerProvisionerAPIMockRecorder is the mock recorder for MockContainerProvisionerAPI.
type MockContainerProvisionerAPIMockRecorder struct {
	mock *MockContainerProvisionerAPI
}

// NewMockContainerProvisionerAPI creates a new mock instance.
func NewMockContainerProvisionerAPI(ctrl *gomock.Controller) *MockContainerProvisionerAPI {
	mock := &MockContainerProvisionerAPI{ctrl: ctrl}
	mock.recorder = &MockContainerProvisionerAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContainerProvisionerAPI) EXPECT() *MockContainerProvisionerAPIMockRecorder {
	return m.recorder
}

// ContainerConfig mocks base method.
func (m *MockContainerProvisionerAPI) ContainerConfig() (params.ContainerConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerConfig")
	ret0, _ := ret[0].(params.ContainerConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerConfig indicates an expected call of ContainerConfig.
func (mr *MockContainerProvisionerAPIMockRecorder) ContainerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerConfig", reflect.TypeOf((*MockContainerProvisionerAPI)(nil).ContainerConfig))
}

// ContainerManagerConfig mocks base method.
func (m *MockContainerProvisionerAPI) ContainerManagerConfig(arg0 params.ContainerManagerConfigParams) (params.ContainerManagerConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerManagerConfig", arg0)
	ret0, _ := ret[0].(params.ContainerManagerConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerManagerConfig indicates an expected call of ContainerManagerConfig.
func (mr *MockContainerProvisionerAPIMockRecorder) ContainerManagerConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerManagerConfig", reflect.TypeOf((*MockContainerProvisionerAPI)(nil).ContainerManagerConfig), arg0)
}

// GetContainerProfileInfo mocks base method.
func (m *MockContainerProvisionerAPI) GetContainerProfileInfo(arg0 names.MachineTag) ([]*provisioner.LXDProfileResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContainerProfileInfo", arg0)
	ret0, _ := ret[0].([]*provisioner.LXDProfileResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContainerProfileInfo indicates an expected call of GetContainerProfileInfo.
func (mr *MockContainerProvisionerAPIMockRecorder) GetContainerProfileInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContainerProfileInfo", reflect.TypeOf((*MockContainerProvisionerAPI)(nil).GetContainerProfileInfo), arg0)
}

// HostChangesForContainer mocks base method.
func (m *MockContainerProvisionerAPI) HostChangesForContainer(arg0 names.MachineTag) ([]network0.DeviceToBridge, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HostChangesForContainer", arg0)
	ret0, _ := ret[0].([]network0.DeviceToBridge)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// HostChangesForContainer indicates an expected call of HostChangesForContainer.
func (mr *MockContainerProvisionerAPIMockRecorder) HostChangesForContainer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HostChangesForContainer", reflect.TypeOf((*MockContainerProvisionerAPI)(nil).HostChangesForContainer), arg0)
}

// PrepareContainerInterfaceInfo mocks base method.
func (m *MockContainerProvisionerAPI) PrepareContainerInterfaceInfo(arg0 names.MachineTag) (network.InterfaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareContainerInterfaceInfo", arg0)
	ret0, _ := ret[0].(network.InterfaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareContainerInterfaceInfo indicates an expected call of PrepareContainerInterfaceInfo.
func (mr *MockContainerProvisionerAPIMockRecorder) PrepareContainerInterfaceInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareContainerInterfaceInfo", reflect.TypeOf((*MockContainerProvisionerAPI)(nil).PrepareContainerInterfaceInfo), arg0)
}

// ReleaseContainerAddresses mocks base method.
func (m *MockContainerProvisionerAPI) ReleaseContainerAddresses(arg0 names.MachineTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseContainerAddresses", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseContainerAddresses indicates an expected call of ReleaseContainerAddresses.
func (mr *MockContainerProvisionerAPIMockRecorder) ReleaseContainerAddresses(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseContainerAddresses", reflect.TypeOf((*MockContainerProvisionerAPI)(nil).ReleaseContainerAddresses), arg0)
}

// SetHostMachineNetworkConfig mocks base method.
func (m *MockContainerProvisionerAPI) SetHostMachineNetworkConfig(arg0 names.MachineTag, arg1 []params.NetworkConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHostMachineNetworkConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHostMachineNetworkConfig indicates an expected call of SetHostMachineNetworkConfig.
func (mr *MockContainerProvisionerAPIMockRecorder) SetHostMachineNetworkConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHostMachineNetworkConfig", reflect.TypeOf((*MockContainerProvisionerAPI)(nil).SetHostMachineNetworkConfig), arg0, arg1)
}

// MockControllerAPI is a mock of ControllerAPI interface.
type MockControllerAPI struct {
	ctrl     *gomock.Controller
	recorder *MockControllerAPIMockRecorder
}

// MockControllerAPIMockRecorder is the mock recorder for MockControllerAPI.
type MockControllerAPIMockRecorder struct {
	mock *MockControllerAPI
}

// NewMockControllerAPI creates a new mock instance.
func NewMockControllerAPI(ctrl *gomock.Controller) *MockControllerAPI {
	mock := &MockControllerAPI{ctrl: ctrl}
	mock.recorder = &MockControllerAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockControllerAPI) EXPECT() *MockControllerAPIMockRecorder {
	return m.recorder
}

// APIAddresses mocks base method.
func (m *MockControllerAPI) APIAddresses() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIAddresses")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// APIAddresses indicates an expected call of APIAddresses.
func (mr *MockControllerAPIMockRecorder) APIAddresses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIAddresses", reflect.TypeOf((*MockControllerAPI)(nil).APIAddresses))
}

// CACert mocks base method.
func (m *MockControllerAPI) CACert() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CACert")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CACert indicates an expected call of CACert.
func (mr *MockControllerAPIMockRecorder) CACert() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CACert", reflect.TypeOf((*MockControllerAPI)(nil).CACert))
}

// ControllerConfig mocks base method.
func (m *MockControllerAPI) ControllerConfig() (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig")
	ret0, _ := ret[0].(controller.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig.
func (mr *MockControllerAPIMockRecorder) ControllerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockControllerAPI)(nil).ControllerConfig))
}

// ModelConfig mocks base method.
func (m *MockControllerAPI) ModelConfig(arg0 context.Context) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig", arg0)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig.
func (mr *MockControllerAPIMockRecorder) ModelConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockControllerAPI)(nil).ModelConfig), arg0)
}

// ModelUUID mocks base method.
func (m *MockControllerAPI) ModelUUID() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelUUID")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelUUID indicates an expected call of ModelUUID.
func (mr *MockControllerAPIMockRecorder) ModelUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelUUID", reflect.TypeOf((*MockControllerAPI)(nil).ModelUUID))
}

// WatchForModelConfigChanges mocks base method.
func (m *MockControllerAPI) WatchForModelConfigChanges() (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchForModelConfigChanges")
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchForModelConfigChanges indicates an expected call of WatchForModelConfigChanges.
func (mr *MockControllerAPIMockRecorder) WatchForModelConfigChanges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchForModelConfigChanges", reflect.TypeOf((*MockControllerAPI)(nil).WatchForModelConfigChanges))
}

// MockMachinesAPI is a mock of MachinesAPI interface.
type MockMachinesAPI struct {
	ctrl     *gomock.Controller
	recorder *MockMachinesAPIMockRecorder
}

// MockMachinesAPIMockRecorder is the mock recorder for MockMachinesAPI.
type MockMachinesAPIMockRecorder struct {
	mock *MockMachinesAPI
}

// NewMockMachinesAPI creates a new mock instance.
func NewMockMachinesAPI(ctrl *gomock.Controller) *MockMachinesAPI {
	mock := &MockMachinesAPI{ctrl: ctrl}
	mock.recorder = &MockMachinesAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachinesAPI) EXPECT() *MockMachinesAPIMockRecorder {
	return m.recorder
}

// Machines mocks base method.
func (m *MockMachinesAPI) Machines(arg0 ...names.MachineTag) ([]provisioner.MachineResult, error) {
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
func (mr *MockMachinesAPIMockRecorder) Machines(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machines", reflect.TypeOf((*MockMachinesAPI)(nil).Machines), arg0...)
}

// MachinesWithTransientErrors mocks base method.
func (m *MockMachinesAPI) MachinesWithTransientErrors() ([]provisioner.MachineStatusResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MachinesWithTransientErrors")
	ret0, _ := ret[0].([]provisioner.MachineStatusResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MachinesWithTransientErrors indicates an expected call of MachinesWithTransientErrors.
func (mr *MockMachinesAPIMockRecorder) MachinesWithTransientErrors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MachinesWithTransientErrors", reflect.TypeOf((*MockMachinesAPI)(nil).MachinesWithTransientErrors))
}

// ProvisioningInfo mocks base method.
func (m *MockMachinesAPI) ProvisioningInfo(arg0 []names.MachineTag) (params.ProvisioningInfoResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProvisioningInfo", arg0)
	ret0, _ := ret[0].(params.ProvisioningInfoResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProvisioningInfo indicates an expected call of ProvisioningInfo.
func (mr *MockMachinesAPIMockRecorder) ProvisioningInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProvisioningInfo", reflect.TypeOf((*MockMachinesAPI)(nil).ProvisioningInfo), arg0)
}

// WatchMachineErrorRetry mocks base method.
func (m *MockMachinesAPI) WatchMachineErrorRetry() (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchMachineErrorRetry")
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchMachineErrorRetry indicates an expected call of WatchMachineErrorRetry.
func (mr *MockMachinesAPIMockRecorder) WatchMachineErrorRetry() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchMachineErrorRetry", reflect.TypeOf((*MockMachinesAPI)(nil).WatchMachineErrorRetry))
}

// WatchModelMachines mocks base method.
func (m *MockMachinesAPI) WatchModelMachines() (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchModelMachines")
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchModelMachines indicates an expected call of WatchModelMachines.
func (mr *MockMachinesAPIMockRecorder) WatchModelMachines() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchModelMachines", reflect.TypeOf((*MockMachinesAPI)(nil).WatchModelMachines))
}
