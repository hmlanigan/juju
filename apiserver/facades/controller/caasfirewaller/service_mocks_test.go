// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/controller/caasfirewaller (interfaces: ApplicationService,PortService)
//
// Generated by this command:
//
//	mockgen -typed -package caasfirewaller_test -destination service_mocks_test.go github.com/juju/juju/apiserver/facades/controller/caasfirewaller ApplicationService,PortService
//

// Package caasfirewaller_test is a generated GoMock package.
package caasfirewaller_test

import (
	context "context"
	reflect "reflect"

	application "github.com/juju/juju/core/application"
	life "github.com/juju/juju/core/life"
	network "github.com/juju/juju/core/network"
	watcher "github.com/juju/juju/core/watcher"
	gomock "go.uber.org/mock/gomock"
)

// MockApplicationService is a mock of ApplicationService interface.
type MockApplicationService struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationServiceMockRecorder
}

// MockApplicationServiceMockRecorder is the mock recorder for MockApplicationService.
type MockApplicationServiceMockRecorder struct {
	mock *MockApplicationService
}

// NewMockApplicationService creates a new mock instance.
func NewMockApplicationService(ctrl *gomock.Controller) *MockApplicationService {
	mock := &MockApplicationService{ctrl: ctrl}
	mock.recorder = &MockApplicationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplicationService) EXPECT() *MockApplicationServiceMockRecorder {
	return m.recorder
}

// GetApplicationIDByName mocks base method.
func (m *MockApplicationService) GetApplicationIDByName(arg0 context.Context, arg1 string) (application.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationIDByName", arg0, arg1)
	ret0, _ := ret[0].(application.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationIDByName indicates an expected call of GetApplicationIDByName.
func (mr *MockApplicationServiceMockRecorder) GetApplicationIDByName(arg0, arg1 any) *MockApplicationServiceGetApplicationIDByNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationIDByName", reflect.TypeOf((*MockApplicationService)(nil).GetApplicationIDByName), arg0, arg1)
	return &MockApplicationServiceGetApplicationIDByNameCall{Call: call}
}

// MockApplicationServiceGetApplicationIDByNameCall wrap *gomock.Call
type MockApplicationServiceGetApplicationIDByNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockApplicationServiceGetApplicationIDByNameCall) Return(arg0 application.ID, arg1 error) *MockApplicationServiceGetApplicationIDByNameCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockApplicationServiceGetApplicationIDByNameCall) Do(f func(context.Context, string) (application.ID, error)) *MockApplicationServiceGetApplicationIDByNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockApplicationServiceGetApplicationIDByNameCall) DoAndReturn(f func(context.Context, string) (application.ID, error)) *MockApplicationServiceGetApplicationIDByNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetApplicationLife mocks base method.
func (m *MockApplicationService) GetApplicationLife(arg0 context.Context, arg1 string) (life.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationLife", arg0, arg1)
	ret0, _ := ret[0].(life.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationLife indicates an expected call of GetApplicationLife.
func (mr *MockApplicationServiceMockRecorder) GetApplicationLife(arg0, arg1 any) *MockApplicationServiceGetApplicationLifeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationLife", reflect.TypeOf((*MockApplicationService)(nil).GetApplicationLife), arg0, arg1)
	return &MockApplicationServiceGetApplicationLifeCall{Call: call}
}

// MockApplicationServiceGetApplicationLifeCall wrap *gomock.Call
type MockApplicationServiceGetApplicationLifeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockApplicationServiceGetApplicationLifeCall) Return(arg0 life.Value, arg1 error) *MockApplicationServiceGetApplicationLifeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockApplicationServiceGetApplicationLifeCall) Do(f func(context.Context, string) (life.Value, error)) *MockApplicationServiceGetApplicationLifeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockApplicationServiceGetApplicationLifeCall) DoAndReturn(f func(context.Context, string) (life.Value, error)) *MockApplicationServiceGetApplicationLifeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetUnitLife mocks base method.
func (m *MockApplicationService) GetUnitLife(arg0 context.Context, arg1 string) (life.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitLife", arg0, arg1)
	ret0, _ := ret[0].(life.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitLife indicates an expected call of GetUnitLife.
func (mr *MockApplicationServiceMockRecorder) GetUnitLife(arg0, arg1 any) *MockApplicationServiceGetUnitLifeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitLife", reflect.TypeOf((*MockApplicationService)(nil).GetUnitLife), arg0, arg1)
	return &MockApplicationServiceGetUnitLifeCall{Call: call}
}

// MockApplicationServiceGetUnitLifeCall wrap *gomock.Call
type MockApplicationServiceGetUnitLifeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockApplicationServiceGetUnitLifeCall) Return(arg0 life.Value, arg1 error) *MockApplicationServiceGetUnitLifeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockApplicationServiceGetUnitLifeCall) Do(f func(context.Context, string) (life.Value, error)) *MockApplicationServiceGetUnitLifeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockApplicationServiceGetUnitLifeCall) DoAndReturn(f func(context.Context, string) (life.Value, error)) *MockApplicationServiceGetUnitLifeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockPortService is a mock of PortService interface.
type MockPortService struct {
	ctrl     *gomock.Controller
	recorder *MockPortServiceMockRecorder
}

// MockPortServiceMockRecorder is the mock recorder for MockPortService.
type MockPortServiceMockRecorder struct {
	mock *MockPortService
}

// NewMockPortService creates a new mock instance.
func NewMockPortService(ctrl *gomock.Controller) *MockPortService {
	mock := &MockPortService{ctrl: ctrl}
	mock.recorder = &MockPortServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPortService) EXPECT() *MockPortServiceMockRecorder {
	return m.recorder
}

// GetApplicationOpenedPortsByEndpoint mocks base method.
func (m *MockPortService) GetApplicationOpenedPortsByEndpoint(arg0 context.Context, arg1 application.ID) (network.GroupedPortRanges, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationOpenedPortsByEndpoint", arg0, arg1)
	ret0, _ := ret[0].(network.GroupedPortRanges)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationOpenedPortsByEndpoint indicates an expected call of GetApplicationOpenedPortsByEndpoint.
func (mr *MockPortServiceMockRecorder) GetApplicationOpenedPortsByEndpoint(arg0, arg1 any) *MockPortServiceGetApplicationOpenedPortsByEndpointCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationOpenedPortsByEndpoint", reflect.TypeOf((*MockPortService)(nil).GetApplicationOpenedPortsByEndpoint), arg0, arg1)
	return &MockPortServiceGetApplicationOpenedPortsByEndpointCall{Call: call}
}

// MockPortServiceGetApplicationOpenedPortsByEndpointCall wrap *gomock.Call
type MockPortServiceGetApplicationOpenedPortsByEndpointCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockPortServiceGetApplicationOpenedPortsByEndpointCall) Return(arg0 network.GroupedPortRanges, arg1 error) *MockPortServiceGetApplicationOpenedPortsByEndpointCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockPortServiceGetApplicationOpenedPortsByEndpointCall) Do(f func(context.Context, application.ID) (network.GroupedPortRanges, error)) *MockPortServiceGetApplicationOpenedPortsByEndpointCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockPortServiceGetApplicationOpenedPortsByEndpointCall) DoAndReturn(f func(context.Context, application.ID) (network.GroupedPortRanges, error)) *MockPortServiceGetApplicationOpenedPortsByEndpointCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WatchApplicationOpenedPorts mocks base method.
func (m *MockPortService) WatchApplicationOpenedPorts(arg0 context.Context) (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchApplicationOpenedPorts", arg0)
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchApplicationOpenedPorts indicates an expected call of WatchApplicationOpenedPorts.
func (mr *MockPortServiceMockRecorder) WatchApplicationOpenedPorts(arg0 any) *MockPortServiceWatchApplicationOpenedPortsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchApplicationOpenedPorts", reflect.TypeOf((*MockPortService)(nil).WatchApplicationOpenedPorts), arg0)
	return &MockPortServiceWatchApplicationOpenedPortsCall{Call: call}
}

// MockPortServiceWatchApplicationOpenedPortsCall wrap *gomock.Call
type MockPortServiceWatchApplicationOpenedPortsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockPortServiceWatchApplicationOpenedPortsCall) Return(arg0 watcher.Watcher[[]string], arg1 error) *MockPortServiceWatchApplicationOpenedPortsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockPortServiceWatchApplicationOpenedPortsCall) Do(f func(context.Context) (watcher.Watcher[[]string], error)) *MockPortServiceWatchApplicationOpenedPortsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockPortServiceWatchApplicationOpenedPortsCall) DoAndReturn(f func(context.Context) (watcher.Watcher[[]string], error)) *MockPortServiceWatchApplicationOpenedPortsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
