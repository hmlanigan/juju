// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/network/modelmigration (interfaces: Coordinator,ImportService,ExportService)
//
// Generated by this command:
//
//	mockgen -package modelmigration -destination migrations_mock_test.go github.com/juju/juju/domain/network/modelmigration Coordinator,ImportService,ExportService
//

// Package modelmigration is a generated GoMock package.
package modelmigration

import (
	context "context"
	reflect "reflect"

	modelmigration "github.com/juju/juju/core/modelmigration"
	network "github.com/juju/juju/core/network"
	gomock "go.uber.org/mock/gomock"
)

// MockCoordinator is a mock of Coordinator interface.
type MockCoordinator struct {
	ctrl     *gomock.Controller
	recorder *MockCoordinatorMockRecorder
}

// MockCoordinatorMockRecorder is the mock recorder for MockCoordinator.
type MockCoordinatorMockRecorder struct {
	mock *MockCoordinator
}

// NewMockCoordinator creates a new mock instance.
func NewMockCoordinator(ctrl *gomock.Controller) *MockCoordinator {
	mock := &MockCoordinator{ctrl: ctrl}
	mock.recorder = &MockCoordinatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoordinator) EXPECT() *MockCoordinatorMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCoordinator) Add(arg0 modelmigration.Operation) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add.
func (mr *MockCoordinatorMockRecorder) Add(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCoordinator)(nil).Add), arg0)
}

// MockImportService is a mock of ImportService interface.
type MockImportService struct {
	ctrl     *gomock.Controller
	recorder *MockImportServiceMockRecorder
}

// MockImportServiceMockRecorder is the mock recorder for MockImportService.
type MockImportServiceMockRecorder struct {
	mock *MockImportService
}

// NewMockImportService creates a new mock instance.
func NewMockImportService(ctrl *gomock.Controller) *MockImportService {
	mock := &MockImportService{ctrl: ctrl}
	mock.recorder = &MockImportServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImportService) EXPECT() *MockImportServiceMockRecorder {
	return m.recorder
}

// AddSpace mocks base method.
func (m *MockImportService) AddSpace(arg0 context.Context, arg1 network.SpaceInfo) (network.Id, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSpace", arg0, arg1)
	ret0, _ := ret[0].(network.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSpace indicates an expected call of AddSpace.
func (mr *MockImportServiceMockRecorder) AddSpace(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSpace", reflect.TypeOf((*MockImportService)(nil).AddSpace), arg0, arg1)
}

// AddSubnet mocks base method.
func (m *MockImportService) AddSubnet(arg0 context.Context, arg1 network.SubnetInfo) (network.Id, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSubnet", arg0, arg1)
	ret0, _ := ret[0].(network.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSubnet indicates an expected call of AddSubnet.
func (mr *MockImportServiceMockRecorder) AddSubnet(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubnet", reflect.TypeOf((*MockImportService)(nil).AddSubnet), arg0, arg1)
}

// Space mocks base method.
func (m *MockImportService) Space(arg0 context.Context, arg1 string) (*network.SpaceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Space", arg0, arg1)
	ret0, _ := ret[0].(*network.SpaceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Space indicates an expected call of Space.
func (mr *MockImportServiceMockRecorder) Space(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Space", reflect.TypeOf((*MockImportService)(nil).Space), arg0, arg1)
}

// MockExportService is a mock of ExportService interface.
type MockExportService struct {
	ctrl     *gomock.Controller
	recorder *MockExportServiceMockRecorder
}

// MockExportServiceMockRecorder is the mock recorder for MockExportService.
type MockExportServiceMockRecorder struct {
	mock *MockExportService
}

// NewMockExportService creates a new mock instance.
func NewMockExportService(ctrl *gomock.Controller) *MockExportService {
	mock := &MockExportService{ctrl: ctrl}
	mock.recorder = &MockExportServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExportService) EXPECT() *MockExportServiceMockRecorder {
	return m.recorder
}

// GetAllSpaces mocks base method.
func (m *MockExportService) GetAllSpaces(arg0 context.Context) (network.SpaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSpaces", arg0)
	ret0, _ := ret[0].(network.SpaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSpaces indicates an expected call of GetAllSpaces.
func (mr *MockExportServiceMockRecorder) GetAllSpaces(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSpaces", reflect.TypeOf((*MockExportService)(nil).GetAllSpaces), arg0)
}

// GetAllSubnets mocks base method.
func (m *MockExportService) GetAllSubnets(arg0 context.Context) (network.SubnetInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSubnets", arg0)
	ret0, _ := ret[0].(network.SubnetInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSubnets indicates an expected call of GetAllSubnets.
func (mr *MockExportServiceMockRecorder) GetAllSubnets(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSubnets", reflect.TypeOf((*MockExportService)(nil).GetAllSubnets), arg0)
}