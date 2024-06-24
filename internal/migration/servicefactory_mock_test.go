// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/servicefactory (interfaces: ServiceFactoryGetter,ServiceFactory)
//
// Generated by this command:
//
//	mockgen -typed -package migration_test -destination servicefactory_mock_test.go github.com/juju/juju/internal/servicefactory ServiceFactoryGetter,ServiceFactory
//

// Package migration_test is a generated GoMock package.
package migration_test

import (
	reflect "reflect"

	service "github.com/juju/juju/domain/access/service"
	service0 "github.com/juju/juju/domain/annotation/service"
	service1 "github.com/juju/juju/domain/application/service"
	service2 "github.com/juju/juju/domain/autocert/service"
	service3 "github.com/juju/juju/domain/bakerystorage/service"
	service4 "github.com/juju/juju/domain/blockdevice/service"
	service5 "github.com/juju/juju/domain/cloud/service"
	service6 "github.com/juju/juju/domain/controllerconfig/service"
	service7 "github.com/juju/juju/domain/controllernode/service"
	service8 "github.com/juju/juju/domain/credential/service"
	service9 "github.com/juju/juju/domain/externalcontroller/service"
	service10 "github.com/juju/juju/domain/flag/service"
	service11 "github.com/juju/juju/domain/machine/service"
	service12 "github.com/juju/juju/domain/model/service"
	service13 "github.com/juju/juju/domain/modelconfig/service"
	service14 "github.com/juju/juju/domain/modeldefaults/service"
	service15 "github.com/juju/juju/domain/network/service"
	service16 "github.com/juju/juju/domain/objectstore/service"
	service17 "github.com/juju/juju/domain/secret/service"
	service18 "github.com/juju/juju/domain/secretbackend/service"
	service19 "github.com/juju/juju/domain/storage/service"
	service20 "github.com/juju/juju/domain/unit/service"
	service21 "github.com/juju/juju/domain/upgrade/service"
	servicefactory "github.com/juju/juju/internal/servicefactory"
	storage "github.com/juju/juju/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockServiceFactoryGetter is a mock of ServiceFactoryGetter interface.
type MockServiceFactoryGetter struct {
	ctrl     *gomock.Controller
	recorder *MockServiceFactoryGetterMockRecorder
}

// MockServiceFactoryGetterMockRecorder is the mock recorder for MockServiceFactoryGetter.
type MockServiceFactoryGetterMockRecorder struct {
	mock *MockServiceFactoryGetter
}

// NewMockServiceFactoryGetter creates a new mock instance.
func NewMockServiceFactoryGetter(ctrl *gomock.Controller) *MockServiceFactoryGetter {
	mock := &MockServiceFactoryGetter{ctrl: ctrl}
	mock.recorder = &MockServiceFactoryGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceFactoryGetter) EXPECT() *MockServiceFactoryGetterMockRecorder {
	return m.recorder
}

// FactoryForModel mocks base method.
func (m *MockServiceFactoryGetter) FactoryForModel(arg0 string) servicefactory.ServiceFactory {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FactoryForModel", arg0)
	ret0, _ := ret[0].(servicefactory.ServiceFactory)
	return ret0
}

// FactoryForModel indicates an expected call of FactoryForModel.
func (mr *MockServiceFactoryGetterMockRecorder) FactoryForModel(arg0 any) *MockServiceFactoryGetterFactoryForModelCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FactoryForModel", reflect.TypeOf((*MockServiceFactoryGetter)(nil).FactoryForModel), arg0)
	return &MockServiceFactoryGetterFactoryForModelCall{Call: call}
}

// MockServiceFactoryGetterFactoryForModelCall wrap *gomock.Call
type MockServiceFactoryGetterFactoryForModelCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryGetterFactoryForModelCall) Return(arg0 servicefactory.ServiceFactory) *MockServiceFactoryGetterFactoryForModelCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryGetterFactoryForModelCall) Do(f func(string) servicefactory.ServiceFactory) *MockServiceFactoryGetterFactoryForModelCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryGetterFactoryForModelCall) DoAndReturn(f func(string) servicefactory.ServiceFactory) *MockServiceFactoryGetterFactoryForModelCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockServiceFactory is a mock of ServiceFactory interface.
type MockServiceFactory struct {
	ctrl     *gomock.Controller
	recorder *MockServiceFactoryMockRecorder
}

// MockServiceFactoryMockRecorder is the mock recorder for MockServiceFactory.
type MockServiceFactoryMockRecorder struct {
	mock *MockServiceFactory
}

// NewMockServiceFactory creates a new mock instance.
func NewMockServiceFactory(ctrl *gomock.Controller) *MockServiceFactory {
	mock := &MockServiceFactory{ctrl: ctrl}
	mock.recorder = &MockServiceFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceFactory) EXPECT() *MockServiceFactoryMockRecorder {
	return m.recorder
}

// Access mocks base method.
func (m *MockServiceFactory) Access() *service.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Access")
	ret0, _ := ret[0].(*service.Service)
	return ret0
}

// Access indicates an expected call of Access.
func (mr *MockServiceFactoryMockRecorder) Access() *MockServiceFactoryAccessCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Access", reflect.TypeOf((*MockServiceFactory)(nil).Access))
	return &MockServiceFactoryAccessCall{Call: call}
}

// MockServiceFactoryAccessCall wrap *gomock.Call
type MockServiceFactoryAccessCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryAccessCall) Return(arg0 *service.Service) *MockServiceFactoryAccessCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryAccessCall) Do(f func() *service.Service) *MockServiceFactoryAccessCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryAccessCall) DoAndReturn(f func() *service.Service) *MockServiceFactoryAccessCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// AgentObjectStore mocks base method.
func (m *MockServiceFactory) AgentObjectStore() *service16.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentObjectStore")
	ret0, _ := ret[0].(*service16.WatchableService)
	return ret0
}

// AgentObjectStore indicates an expected call of AgentObjectStore.
func (mr *MockServiceFactoryMockRecorder) AgentObjectStore() *MockServiceFactoryAgentObjectStoreCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentObjectStore", reflect.TypeOf((*MockServiceFactory)(nil).AgentObjectStore))
	return &MockServiceFactoryAgentObjectStoreCall{Call: call}
}

// MockServiceFactoryAgentObjectStoreCall wrap *gomock.Call
type MockServiceFactoryAgentObjectStoreCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryAgentObjectStoreCall) Return(arg0 *service16.WatchableService) *MockServiceFactoryAgentObjectStoreCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryAgentObjectStoreCall) Do(f func() *service16.WatchableService) *MockServiceFactoryAgentObjectStoreCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryAgentObjectStoreCall) DoAndReturn(f func() *service16.WatchableService) *MockServiceFactoryAgentObjectStoreCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Annotation mocks base method.
func (m *MockServiceFactory) Annotation() *service0.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Annotation")
	ret0, _ := ret[0].(*service0.Service)
	return ret0
}

// Annotation indicates an expected call of Annotation.
func (mr *MockServiceFactoryMockRecorder) Annotation() *MockServiceFactoryAnnotationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Annotation", reflect.TypeOf((*MockServiceFactory)(nil).Annotation))
	return &MockServiceFactoryAnnotationCall{Call: call}
}

// MockServiceFactoryAnnotationCall wrap *gomock.Call
type MockServiceFactoryAnnotationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryAnnotationCall) Return(arg0 *service0.Service) *MockServiceFactoryAnnotationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryAnnotationCall) Do(f func() *service0.Service) *MockServiceFactoryAnnotationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryAnnotationCall) DoAndReturn(f func() *service0.Service) *MockServiceFactoryAnnotationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Application mocks base method.
func (m *MockServiceFactory) Application(arg0 storage.ProviderRegistry) *service1.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application", arg0)
	ret0, _ := ret[0].(*service1.Service)
	return ret0
}

// Application indicates an expected call of Application.
func (mr *MockServiceFactoryMockRecorder) Application(arg0 any) *MockServiceFactoryApplicationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockServiceFactory)(nil).Application), arg0)
	return &MockServiceFactoryApplicationCall{Call: call}
}

// MockServiceFactoryApplicationCall wrap *gomock.Call
type MockServiceFactoryApplicationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryApplicationCall) Return(arg0 *service1.Service) *MockServiceFactoryApplicationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryApplicationCall) Do(f func(storage.ProviderRegistry) *service1.Service) *MockServiceFactoryApplicationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryApplicationCall) DoAndReturn(f func(storage.ProviderRegistry) *service1.Service) *MockServiceFactoryApplicationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// AutocertCache mocks base method.
func (m *MockServiceFactory) AutocertCache() *service2.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AutocertCache")
	ret0, _ := ret[0].(*service2.Service)
	return ret0
}

// AutocertCache indicates an expected call of AutocertCache.
func (mr *MockServiceFactoryMockRecorder) AutocertCache() *MockServiceFactoryAutocertCacheCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutocertCache", reflect.TypeOf((*MockServiceFactory)(nil).AutocertCache))
	return &MockServiceFactoryAutocertCacheCall{Call: call}
}

// MockServiceFactoryAutocertCacheCall wrap *gomock.Call
type MockServiceFactoryAutocertCacheCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryAutocertCacheCall) Return(arg0 *service2.Service) *MockServiceFactoryAutocertCacheCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryAutocertCacheCall) Do(f func() *service2.Service) *MockServiceFactoryAutocertCacheCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryAutocertCacheCall) DoAndReturn(f func() *service2.Service) *MockServiceFactoryAutocertCacheCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BlockDevice mocks base method.
func (m *MockServiceFactory) BlockDevice() *service4.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockDevice")
	ret0, _ := ret[0].(*service4.WatchableService)
	return ret0
}

// BlockDevice indicates an expected call of BlockDevice.
func (mr *MockServiceFactoryMockRecorder) BlockDevice() *MockServiceFactoryBlockDeviceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockDevice", reflect.TypeOf((*MockServiceFactory)(nil).BlockDevice))
	return &MockServiceFactoryBlockDeviceCall{Call: call}
}

// MockServiceFactoryBlockDeviceCall wrap *gomock.Call
type MockServiceFactoryBlockDeviceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryBlockDeviceCall) Return(arg0 *service4.WatchableService) *MockServiceFactoryBlockDeviceCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryBlockDeviceCall) Do(f func() *service4.WatchableService) *MockServiceFactoryBlockDeviceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryBlockDeviceCall) DoAndReturn(f func() *service4.WatchableService) *MockServiceFactoryBlockDeviceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Cloud mocks base method.
func (m *MockServiceFactory) Cloud() *service5.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cloud")
	ret0, _ := ret[0].(*service5.WatchableService)
	return ret0
}

// Cloud indicates an expected call of Cloud.
func (mr *MockServiceFactoryMockRecorder) Cloud() *MockServiceFactoryCloudCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cloud", reflect.TypeOf((*MockServiceFactory)(nil).Cloud))
	return &MockServiceFactoryCloudCall{Call: call}
}

// MockServiceFactoryCloudCall wrap *gomock.Call
type MockServiceFactoryCloudCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryCloudCall) Return(arg0 *service5.WatchableService) *MockServiceFactoryCloudCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryCloudCall) Do(f func() *service5.WatchableService) *MockServiceFactoryCloudCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryCloudCall) DoAndReturn(f func() *service5.WatchableService) *MockServiceFactoryCloudCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Config mocks base method.
func (m *MockServiceFactory) Config() *service13.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*service13.WatchableService)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockServiceFactoryMockRecorder) Config() *MockServiceFactoryConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockServiceFactory)(nil).Config))
	return &MockServiceFactoryConfigCall{Call: call}
}

// MockServiceFactoryConfigCall wrap *gomock.Call
type MockServiceFactoryConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryConfigCall) Return(arg0 *service13.WatchableService) *MockServiceFactoryConfigCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryConfigCall) Do(f func() *service13.WatchableService) *MockServiceFactoryConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryConfigCall) DoAndReturn(f func() *service13.WatchableService) *MockServiceFactoryConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ControllerConfig mocks base method.
func (m *MockServiceFactory) ControllerConfig() *service6.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig")
	ret0, _ := ret[0].(*service6.WatchableService)
	return ret0
}

// ControllerConfig indicates an expected call of ControllerConfig.
func (mr *MockServiceFactoryMockRecorder) ControllerConfig() *MockServiceFactoryControllerConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockServiceFactory)(nil).ControllerConfig))
	return &MockServiceFactoryControllerConfigCall{Call: call}
}

// MockServiceFactoryControllerConfigCall wrap *gomock.Call
type MockServiceFactoryControllerConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryControllerConfigCall) Return(arg0 *service6.WatchableService) *MockServiceFactoryControllerConfigCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryControllerConfigCall) Do(f func() *service6.WatchableService) *MockServiceFactoryControllerConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryControllerConfigCall) DoAndReturn(f func() *service6.WatchableService) *MockServiceFactoryControllerConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ControllerNode mocks base method.
func (m *MockServiceFactory) ControllerNode() *service7.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerNode")
	ret0, _ := ret[0].(*service7.Service)
	return ret0
}

// ControllerNode indicates an expected call of ControllerNode.
func (mr *MockServiceFactoryMockRecorder) ControllerNode() *MockServiceFactoryControllerNodeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerNode", reflect.TypeOf((*MockServiceFactory)(nil).ControllerNode))
	return &MockServiceFactoryControllerNodeCall{Call: call}
}

// MockServiceFactoryControllerNodeCall wrap *gomock.Call
type MockServiceFactoryControllerNodeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryControllerNodeCall) Return(arg0 *service7.Service) *MockServiceFactoryControllerNodeCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryControllerNodeCall) Do(f func() *service7.Service) *MockServiceFactoryControllerNodeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryControllerNodeCall) DoAndReturn(f func() *service7.Service) *MockServiceFactoryControllerNodeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Credential mocks base method.
func (m *MockServiceFactory) Credential() *service8.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Credential")
	ret0, _ := ret[0].(*service8.WatchableService)
	return ret0
}

// Credential indicates an expected call of Credential.
func (mr *MockServiceFactoryMockRecorder) Credential() *MockServiceFactoryCredentialCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Credential", reflect.TypeOf((*MockServiceFactory)(nil).Credential))
	return &MockServiceFactoryCredentialCall{Call: call}
}

// MockServiceFactoryCredentialCall wrap *gomock.Call
type MockServiceFactoryCredentialCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryCredentialCall) Return(arg0 *service8.WatchableService) *MockServiceFactoryCredentialCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryCredentialCall) Do(f func() *service8.WatchableService) *MockServiceFactoryCredentialCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryCredentialCall) DoAndReturn(f func() *service8.WatchableService) *MockServiceFactoryCredentialCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ExternalController mocks base method.
func (m *MockServiceFactory) ExternalController() *service9.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExternalController")
	ret0, _ := ret[0].(*service9.WatchableService)
	return ret0
}

// ExternalController indicates an expected call of ExternalController.
func (mr *MockServiceFactoryMockRecorder) ExternalController() *MockServiceFactoryExternalControllerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExternalController", reflect.TypeOf((*MockServiceFactory)(nil).ExternalController))
	return &MockServiceFactoryExternalControllerCall{Call: call}
}

// MockServiceFactoryExternalControllerCall wrap *gomock.Call
type MockServiceFactoryExternalControllerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryExternalControllerCall) Return(arg0 *service9.WatchableService) *MockServiceFactoryExternalControllerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryExternalControllerCall) Do(f func() *service9.WatchableService) *MockServiceFactoryExternalControllerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryExternalControllerCall) DoAndReturn(f func() *service9.WatchableService) *MockServiceFactoryExternalControllerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Flag mocks base method.
func (m *MockServiceFactory) Flag() *service10.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flag")
	ret0, _ := ret[0].(*service10.Service)
	return ret0
}

// Flag indicates an expected call of Flag.
func (mr *MockServiceFactoryMockRecorder) Flag() *MockServiceFactoryFlagCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flag", reflect.TypeOf((*MockServiceFactory)(nil).Flag))
	return &MockServiceFactoryFlagCall{Call: call}
}

// MockServiceFactoryFlagCall wrap *gomock.Call
type MockServiceFactoryFlagCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryFlagCall) Return(arg0 *service10.Service) *MockServiceFactoryFlagCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryFlagCall) Do(f func() *service10.Service) *MockServiceFactoryFlagCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryFlagCall) DoAndReturn(f func() *service10.Service) *MockServiceFactoryFlagCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Macaroon mocks base method.
func (m *MockServiceFactory) Macaroon() *service3.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Macaroon")
	ret0, _ := ret[0].(*service3.Service)
	return ret0
}

// Macaroon indicates an expected call of Macaroon.
func (mr *MockServiceFactoryMockRecorder) Macaroon() *MockServiceFactoryMacaroonCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Macaroon", reflect.TypeOf((*MockServiceFactory)(nil).Macaroon))
	return &MockServiceFactoryMacaroonCall{Call: call}
}

// MockServiceFactoryMacaroonCall wrap *gomock.Call
type MockServiceFactoryMacaroonCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryMacaroonCall) Return(arg0 *service3.Service) *MockServiceFactoryMacaroonCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryMacaroonCall) Do(f func() *service3.Service) *MockServiceFactoryMacaroonCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryMacaroonCall) DoAndReturn(f func() *service3.Service) *MockServiceFactoryMacaroonCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Machine mocks base method.
func (m *MockServiceFactory) Machine() *service11.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Machine")
	ret0, _ := ret[0].(*service11.Service)
	return ret0
}

// Machine indicates an expected call of Machine.
func (mr *MockServiceFactoryMockRecorder) Machine() *MockServiceFactoryMachineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machine", reflect.TypeOf((*MockServiceFactory)(nil).Machine))
	return &MockServiceFactoryMachineCall{Call: call}
}

// MockServiceFactoryMachineCall wrap *gomock.Call
type MockServiceFactoryMachineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryMachineCall) Return(arg0 *service11.Service) *MockServiceFactoryMachineCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryMachineCall) Do(f func() *service11.Service) *MockServiceFactoryMachineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryMachineCall) DoAndReturn(f func() *service11.Service) *MockServiceFactoryMachineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Model mocks base method.
func (m *MockServiceFactory) Model() *service12.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Model")
	ret0, _ := ret[0].(*service12.Service)
	return ret0
}

// Model indicates an expected call of Model.
func (mr *MockServiceFactoryMockRecorder) Model() *MockServiceFactoryModelCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Model", reflect.TypeOf((*MockServiceFactory)(nil).Model))
	return &MockServiceFactoryModelCall{Call: call}
}

// MockServiceFactoryModelCall wrap *gomock.Call
type MockServiceFactoryModelCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryModelCall) Return(arg0 *service12.Service) *MockServiceFactoryModelCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryModelCall) Do(f func() *service12.Service) *MockServiceFactoryModelCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryModelCall) DoAndReturn(f func() *service12.Service) *MockServiceFactoryModelCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelDefaults mocks base method.
func (m *MockServiceFactory) ModelDefaults() *service14.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelDefaults")
	ret0, _ := ret[0].(*service14.Service)
	return ret0
}

// ModelDefaults indicates an expected call of ModelDefaults.
func (mr *MockServiceFactoryMockRecorder) ModelDefaults() *MockServiceFactoryModelDefaultsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelDefaults", reflect.TypeOf((*MockServiceFactory)(nil).ModelDefaults))
	return &MockServiceFactoryModelDefaultsCall{Call: call}
}

// MockServiceFactoryModelDefaultsCall wrap *gomock.Call
type MockServiceFactoryModelDefaultsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryModelDefaultsCall) Return(arg0 *service14.Service) *MockServiceFactoryModelDefaultsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryModelDefaultsCall) Do(f func() *service14.Service) *MockServiceFactoryModelDefaultsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryModelDefaultsCall) DoAndReturn(f func() *service14.Service) *MockServiceFactoryModelDefaultsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ModelInfo mocks base method.
func (m *MockServiceFactory) ModelInfo() *service12.ModelService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelInfo")
	ret0, _ := ret[0].(*service12.ModelService)
	return ret0
}

// ModelInfo indicates an expected call of ModelInfo.
func (mr *MockServiceFactoryMockRecorder) ModelInfo() *MockServiceFactoryModelInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelInfo", reflect.TypeOf((*MockServiceFactory)(nil).ModelInfo))
	return &MockServiceFactoryModelInfoCall{Call: call}
}

// MockServiceFactoryModelInfoCall wrap *gomock.Call
type MockServiceFactoryModelInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryModelInfoCall) Return(arg0 *service12.ModelService) *MockServiceFactoryModelInfoCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryModelInfoCall) Do(f func() *service12.ModelService) *MockServiceFactoryModelInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryModelInfoCall) DoAndReturn(f func() *service12.ModelService) *MockServiceFactoryModelInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Network mocks base method.
func (m *MockServiceFactory) Network() *service15.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Network")
	ret0, _ := ret[0].(*service15.WatchableService)
	return ret0
}

// Network indicates an expected call of Network.
func (mr *MockServiceFactoryMockRecorder) Network() *MockServiceFactoryNetworkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Network", reflect.TypeOf((*MockServiceFactory)(nil).Network))
	return &MockServiceFactoryNetworkCall{Call: call}
}

// MockServiceFactoryNetworkCall wrap *gomock.Call
type MockServiceFactoryNetworkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryNetworkCall) Return(arg0 *service15.WatchableService) *MockServiceFactoryNetworkCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryNetworkCall) Do(f func() *service15.WatchableService) *MockServiceFactoryNetworkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryNetworkCall) DoAndReturn(f func() *service15.WatchableService) *MockServiceFactoryNetworkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ObjectStore mocks base method.
func (m *MockServiceFactory) ObjectStore() *service16.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObjectStore")
	ret0, _ := ret[0].(*service16.WatchableService)
	return ret0
}

// ObjectStore indicates an expected call of ObjectStore.
func (mr *MockServiceFactoryMockRecorder) ObjectStore() *MockServiceFactoryObjectStoreCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObjectStore", reflect.TypeOf((*MockServiceFactory)(nil).ObjectStore))
	return &MockServiceFactoryObjectStoreCall{Call: call}
}

// MockServiceFactoryObjectStoreCall wrap *gomock.Call
type MockServiceFactoryObjectStoreCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryObjectStoreCall) Return(arg0 *service16.WatchableService) *MockServiceFactoryObjectStoreCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryObjectStoreCall) Do(f func() *service16.WatchableService) *MockServiceFactoryObjectStoreCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryObjectStoreCall) DoAndReturn(f func() *service16.WatchableService) *MockServiceFactoryObjectStoreCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Secret mocks base method.
func (m *MockServiceFactory) Secret(arg0 service17.BackendAdminConfigGetter) *service17.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Secret", arg0)
	ret0, _ := ret[0].(*service17.WatchableService)
	return ret0
}

// Secret indicates an expected call of Secret.
func (mr *MockServiceFactoryMockRecorder) Secret(arg0 any) *MockServiceFactorySecretCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Secret", reflect.TypeOf((*MockServiceFactory)(nil).Secret), arg0)
	return &MockServiceFactorySecretCall{Call: call}
}

// MockServiceFactorySecretCall wrap *gomock.Call
type MockServiceFactorySecretCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactorySecretCall) Return(arg0 *service17.WatchableService) *MockServiceFactorySecretCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactorySecretCall) Do(f func(service17.BackendAdminConfigGetter) *service17.WatchableService) *MockServiceFactorySecretCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactorySecretCall) DoAndReturn(f func(service17.BackendAdminConfigGetter) *service17.WatchableService) *MockServiceFactorySecretCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SecretBackend mocks base method.
func (m *MockServiceFactory) SecretBackend() *service18.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretBackend")
	ret0, _ := ret[0].(*service18.WatchableService)
	return ret0
}

// SecretBackend indicates an expected call of SecretBackend.
func (mr *MockServiceFactoryMockRecorder) SecretBackend() *MockServiceFactorySecretBackendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretBackend", reflect.TypeOf((*MockServiceFactory)(nil).SecretBackend))
	return &MockServiceFactorySecretBackendCall{Call: call}
}

// MockServiceFactorySecretBackendCall wrap *gomock.Call
type MockServiceFactorySecretBackendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactorySecretBackendCall) Return(arg0 *service18.WatchableService) *MockServiceFactorySecretBackendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactorySecretBackendCall) Do(f func() *service18.WatchableService) *MockServiceFactorySecretBackendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactorySecretBackendCall) DoAndReturn(f func() *service18.WatchableService) *MockServiceFactorySecretBackendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Storage mocks base method.
func (m *MockServiceFactory) Storage(arg0 storage.ProviderRegistry) *service19.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Storage", arg0)
	ret0, _ := ret[0].(*service19.Service)
	return ret0
}

// Storage indicates an expected call of Storage.
func (mr *MockServiceFactoryMockRecorder) Storage(arg0 any) *MockServiceFactoryStorageCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Storage", reflect.TypeOf((*MockServiceFactory)(nil).Storage), arg0)
	return &MockServiceFactoryStorageCall{Call: call}
}

// MockServiceFactoryStorageCall wrap *gomock.Call
type MockServiceFactoryStorageCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryStorageCall) Return(arg0 *service19.Service) *MockServiceFactoryStorageCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryStorageCall) Do(f func(storage.ProviderRegistry) *service19.Service) *MockServiceFactoryStorageCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryStorageCall) DoAndReturn(f func(storage.ProviderRegistry) *service19.Service) *MockServiceFactoryStorageCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Unit mocks base method.
func (m *MockServiceFactory) Unit() *service20.Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unit")
	ret0, _ := ret[0].(*service20.Service)
	return ret0
}

// Unit indicates an expected call of Unit.
func (mr *MockServiceFactoryMockRecorder) Unit() *MockServiceFactoryUnitCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unit", reflect.TypeOf((*MockServiceFactory)(nil).Unit))
	return &MockServiceFactoryUnitCall{Call: call}
}

// MockServiceFactoryUnitCall wrap *gomock.Call
type MockServiceFactoryUnitCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryUnitCall) Return(arg0 *service20.Service) *MockServiceFactoryUnitCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryUnitCall) Do(f func() *service20.Service) *MockServiceFactoryUnitCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryUnitCall) DoAndReturn(f func() *service20.Service) *MockServiceFactoryUnitCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Upgrade mocks base method.
func (m *MockServiceFactory) Upgrade() *service21.WatchableService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upgrade")
	ret0, _ := ret[0].(*service21.WatchableService)
	return ret0
}

// Upgrade indicates an expected call of Upgrade.
func (mr *MockServiceFactoryMockRecorder) Upgrade() *MockServiceFactoryUpgradeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upgrade", reflect.TypeOf((*MockServiceFactory)(nil).Upgrade))
	return &MockServiceFactoryUpgradeCall{Call: call}
}

// MockServiceFactoryUpgradeCall wrap *gomock.Call
type MockServiceFactoryUpgradeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceFactoryUpgradeCall) Return(arg0 *service21.WatchableService) *MockServiceFactoryUpgradeCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceFactoryUpgradeCall) Do(f func() *service21.WatchableService) *MockServiceFactoryUpgradeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceFactoryUpgradeCall) DoAndReturn(f func() *service21.WatchableService) *MockServiceFactoryUpgradeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
