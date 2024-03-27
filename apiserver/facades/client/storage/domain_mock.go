// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/storage (interfaces: StorageService)
//
// Generated by this command:
//
//	mockgen -package storage -destination domain_mock.go github.com/juju/juju/apiserver/facades/client/storage StorageService
//

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	storage "github.com/juju/juju/domain/storage"
	service "github.com/juju/juju/domain/storage/service"
	storage0 "github.com/juju/juju/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockStorageService is a mock of StorageService interface.
type MockStorageService struct {
	ctrl     *gomock.Controller
	recorder *MockStorageServiceMockRecorder
}

// MockStorageServiceMockRecorder is the mock recorder for MockStorageService.
type MockStorageServiceMockRecorder struct {
	mock *MockStorageService
}

// NewMockStorageService creates a new mock instance.
func NewMockStorageService(ctrl *gomock.Controller) *MockStorageService {
	mock := &MockStorageService{ctrl: ctrl}
	mock.recorder = &MockStorageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageService) EXPECT() *MockStorageServiceMockRecorder {
	return m.recorder
}

// CreateStoragePool mocks base method.
func (m *MockStorageService) CreateStoragePool(arg0 context.Context, arg1 string, arg2 storage0.ProviderType, arg3 service.PoolAttrs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStoragePool", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateStoragePool indicates an expected call of CreateStoragePool.
func (mr *MockStorageServiceMockRecorder) CreateStoragePool(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStoragePool", reflect.TypeOf((*MockStorageService)(nil).CreateStoragePool), arg0, arg1, arg2, arg3)
}

// DeleteStoragePool mocks base method.
func (m *MockStorageService) DeleteStoragePool(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStoragePool", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStoragePool indicates an expected call of DeleteStoragePool.
func (mr *MockStorageServiceMockRecorder) DeleteStoragePool(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStoragePool", reflect.TypeOf((*MockStorageService)(nil).DeleteStoragePool), arg0, arg1)
}

// GetStoragePoolByName mocks base method.
func (m *MockStorageService) GetStoragePoolByName(arg0 context.Context, arg1 string) (*storage0.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStoragePoolByName", arg0, arg1)
	ret0, _ := ret[0].(*storage0.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStoragePoolByName indicates an expected call of GetStoragePoolByName.
func (mr *MockStorageServiceMockRecorder) GetStoragePoolByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStoragePoolByName", reflect.TypeOf((*MockStorageService)(nil).GetStoragePoolByName), arg0, arg1)
}

// ListStoragePools mocks base method.
func (m *MockStorageService) ListStoragePools(arg0 context.Context, arg1 storage.Names, arg2 storage.Providers) ([]*storage0.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStoragePools", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*storage0.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStoragePools indicates an expected call of ListStoragePools.
func (mr *MockStorageServiceMockRecorder) ListStoragePools(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStoragePools", reflect.TypeOf((*MockStorageService)(nil).ListStoragePools), arg0, arg1, arg2)
}

// ReplaceStoragePool mocks base method.
func (m *MockStorageService) ReplaceStoragePool(arg0 context.Context, arg1 string, arg2 storage0.ProviderType, arg3 service.PoolAttrs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceStoragePool", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceStoragePool indicates an expected call of ReplaceStoragePool.
func (mr *MockStorageServiceMockRecorder) ReplaceStoragePool(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceStoragePool", reflect.TypeOf((*MockStorageService)(nil).ReplaceStoragePool), arg0, arg1, arg2, arg3)
}