// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/modelmanager (interfaces: AccessService,SecretBackendService,ModelService)
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/service_mock.go github.com/juju/juju/apiserver/facades/client/modelmanager AccessService,SecretBackendService,ModelService
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	credential "github.com/juju/juju/core/credential"
	model "github.com/juju/juju/core/model"
	permission "github.com/juju/juju/core/permission"
	user "github.com/juju/juju/core/user"
	model0 "github.com/juju/juju/domain/model"
	service "github.com/juju/juju/domain/secretbackend/service"
	gomock "go.uber.org/mock/gomock"
)

// MockAccessService is a mock of AccessService interface.
type MockAccessService struct {
	ctrl     *gomock.Controller
	recorder *MockAccessServiceMockRecorder
}

// MockAccessServiceMockRecorder is the mock recorder for MockAccessService.
type MockAccessServiceMockRecorder struct {
	mock *MockAccessService
}

// NewMockAccessService creates a new mock instance.
func NewMockAccessService(ctrl *gomock.Controller) *MockAccessService {
	mock := &MockAccessService{ctrl: ctrl}
	mock.recorder = &MockAccessServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccessService) EXPECT() *MockAccessServiceMockRecorder {
	return m.recorder
}

// GetUserByName mocks base method.
func (m *MockAccessService) GetUserByName(arg0 context.Context, arg1 string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByName", arg0, arg1)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByName indicates an expected call of GetUserByName.
func (mr *MockAccessServiceMockRecorder) GetUserByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByName", reflect.TypeOf((*MockAccessService)(nil).GetUserByName), arg0, arg1)
}

// ReadUserAccessLevelForTarget mocks base method.
func (m *MockAccessService) ReadUserAccessLevelForTarget(arg0 context.Context, arg1 string, arg2 permission.ID) (permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserAccessLevelForTarget", arg0, arg1, arg2)
	ret0, _ := ret[0].(permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserAccessLevelForTarget indicates an expected call of ReadUserAccessLevelForTarget.
func (mr *MockAccessServiceMockRecorder) ReadUserAccessLevelForTarget(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserAccessLevelForTarget", reflect.TypeOf((*MockAccessService)(nil).ReadUserAccessLevelForTarget), arg0, arg1, arg2)
}

// MockSecretBackendService is a mock of SecretBackendService interface.
type MockSecretBackendService struct {
	ctrl     *gomock.Controller
	recorder *MockSecretBackendServiceMockRecorder
}

// MockSecretBackendServiceMockRecorder is the mock recorder for MockSecretBackendService.
type MockSecretBackendServiceMockRecorder struct {
	mock *MockSecretBackendService
}

// NewMockSecretBackendService creates a new mock instance.
func NewMockSecretBackendService(ctrl *gomock.Controller) *MockSecretBackendService {
	mock := &MockSecretBackendService{ctrl: ctrl}
	mock.recorder = &MockSecretBackendServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretBackendService) EXPECT() *MockSecretBackendServiceMockRecorder {
	return m.recorder
}

// BackendSummaryInfo mocks base method.
func (m *MockSecretBackendService) BackendSummaryInfo(arg0 context.Context, arg1, arg2 bool, arg3 ...string) ([]*service.SecretBackendInfo, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BackendSummaryInfo", varargs...)
	ret0, _ := ret[0].([]*service.SecretBackendInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BackendSummaryInfo indicates an expected call of BackendSummaryInfo.
func (mr *MockSecretBackendServiceMockRecorder) BackendSummaryInfo(arg0, arg1, arg2 any, arg3 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BackendSummaryInfo", reflect.TypeOf((*MockSecretBackendService)(nil).BackendSummaryInfo), varargs...)
}

// MockModelService is a mock of ModelService interface.
type MockModelService struct {
	ctrl     *gomock.Controller
	recorder *MockModelServiceMockRecorder
}

// MockModelServiceMockRecorder is the mock recorder for MockModelService.
type MockModelServiceMockRecorder struct {
	mock *MockModelService
}

// NewMockModelService creates a new mock instance.
func NewMockModelService(ctrl *gomock.Controller) *MockModelService {
	mock := &MockModelService{ctrl: ctrl}
	mock.recorder = &MockModelServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelService) EXPECT() *MockModelServiceMockRecorder {
	return m.recorder
}

// CreateModel mocks base method.
func (m *MockModelService) CreateModel(arg0 context.Context, arg1 model0.ModelCreationArgs) (model.UUID, func(context.Context) error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateModel", arg0, arg1)
	ret0, _ := ret[0].(model.UUID)
	ret1, _ := ret[1].(func(context.Context) error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateModel indicates an expected call of CreateModel.
func (mr *MockModelServiceMockRecorder) CreateModel(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateModel", reflect.TypeOf((*MockModelService)(nil).CreateModel), arg0, arg1)
}

// DefaultModelCloudNameAndCredential mocks base method.
func (m *MockModelService) DefaultModelCloudNameAndCredential(arg0 context.Context) (string, credential.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultModelCloudNameAndCredential", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(credential.Key)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DefaultModelCloudNameAndCredential indicates an expected call of DefaultModelCloudNameAndCredential.
func (mr *MockModelServiceMockRecorder) DefaultModelCloudNameAndCredential(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultModelCloudNameAndCredential", reflect.TypeOf((*MockModelService)(nil).DefaultModelCloudNameAndCredential), arg0)
}

// DeleteModel mocks base method.
func (m *MockModelService) DeleteModel(arg0 context.Context, arg1 model.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteModel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteModel indicates an expected call of DeleteModel.
func (mr *MockModelServiceMockRecorder) DeleteModel(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteModel", reflect.TypeOf((*MockModelService)(nil).DeleteModel), arg0, arg1)
}

// ListAllModels mocks base method.
func (m *MockModelService) ListAllModels(arg0 context.Context) ([]model.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllModels", arg0)
	ret0, _ := ret[0].([]model.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllModels indicates an expected call of ListAllModels.
func (mr *MockModelServiceMockRecorder) ListAllModels(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllModels", reflect.TypeOf((*MockModelService)(nil).ListAllModels), arg0)
}

// ListModelsForUser mocks base method.
func (m *MockModelService) ListModelsForUser(arg0 context.Context, arg1 user.UUID) ([]model.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListModelsForUser", arg0, arg1)
	ret0, _ := ret[0].([]model.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListModelsForUser indicates an expected call of ListModelsForUser.
func (mr *MockModelServiceMockRecorder) ListModelsForUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListModelsForUser", reflect.TypeOf((*MockModelService)(nil).ListModelsForUser), arg0, arg1)
}

// ModelType mocks base method.
func (m *MockModelService) ModelType(arg0 context.Context, arg1 model.UUID) (model.ModelType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelType", arg0, arg1)
	ret0, _ := ret[0].(model.ModelType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelType indicates an expected call of ModelType.
func (mr *MockModelServiceMockRecorder) ModelType(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelType", reflect.TypeOf((*MockModelService)(nil).ModelType), arg0, arg1)
}
