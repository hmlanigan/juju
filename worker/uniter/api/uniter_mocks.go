// Code generated by MockGen. Hand edited to fix bad code due to generics.
// Source: github.com/juju/juju/worker/uniter/api (interfaces: UniterClient)

// Package mocks is a generated GoMock package.
package api

import (
	"reflect"
	"time"

	"github.com/juju/charm/v11"
	"github.com/juju/names/v4"
	"go.uber.org/mock/gomock"

	"github.com/juju/juju/api/agent/uniter"
	"github.com/juju/juju/core/application"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/rpc/params"
)

// MockUniterClient is a mock of UniterClient interface.
type MockUniterClient struct {
	ctrl     *gomock.Controller
	recorder *MockUniterClientMockRecorder
}

// MockUniterClientMockRecorder is the mock recorder for MockUniterClient.
type MockUniterClientMockRecorder struct {
	mock *MockUniterClient
}

// NewMockUniterClient creates a new mock instance.
func NewMockUniterClient(ctrl *gomock.Controller) *MockUniterClient {
	mock := &MockUniterClient{ctrl: ctrl}
	mock.recorder = &MockUniterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUniterClient) EXPECT() *MockUniterClientMockRecorder {
	return m.recorder
}

// APIAddresses mocks base method.
func (m *MockUniterClient) APIAddresses() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIAddresses")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// APIAddresses indicates an expected call of APIAddresses.
func (mr *MockUniterClientMockRecorder) APIAddresses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIAddresses", reflect.TypeOf((*MockUniterClient)(nil).APIAddresses))
}

// Action mocks base method.
func (m *MockUniterClient) Action(arg0 names.ActionTag) (*uniter.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Action", arg0)
	ret0, _ := ret[0].(*uniter.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Action indicates an expected call of Action.
func (mr *MockUniterClientMockRecorder) Action(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Action", reflect.TypeOf((*MockUniterClient)(nil).Action), arg0)
}

// ActionBegin mocks base method.
func (m *MockUniterClient) ActionBegin(arg0 names.ActionTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActionBegin", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActionBegin indicates an expected call of ActionBegin.
func (mr *MockUniterClientMockRecorder) ActionBegin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionBegin", reflect.TypeOf((*MockUniterClient)(nil).ActionBegin), arg0)
}

// ActionFinish mocks base method.
func (m *MockUniterClient) ActionFinish(arg0 names.ActionTag, arg1 string, arg2 map[string]interface{}, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActionFinish", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActionFinish indicates an expected call of ActionFinish.
func (mr *MockUniterClientMockRecorder) ActionFinish(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionFinish", reflect.TypeOf((*MockUniterClient)(nil).ActionFinish), arg0, arg1, arg2, arg3)
}

// ActionStatus mocks base method.
func (m *MockUniterClient) ActionStatus(arg0 names.ActionTag) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActionStatus", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ActionStatus indicates an expected call of ActionStatus.
func (mr *MockUniterClientMockRecorder) ActionStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionStatus", reflect.TypeOf((*MockUniterClient)(nil).ActionStatus), arg0)
}

// Application mocks base method.
func (m *MockUniterClient) Application(arg0 names.ApplicationTag) (Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application", arg0)
	ret0, _ := ret[0].(Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Application indicates an expected call of Application.
func (mr *MockUniterClientMockRecorder) Application(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockUniterClient)(nil).Application), arg0)
}

// Charm mocks base method.
func (m *MockUniterClient) Charm(arg0 *charm.URL) (Charm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Charm", arg0)
	ret0, _ := ret[0].(Charm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Charm indicates an expected call of Charm.
func (mr *MockUniterClientMockRecorder) Charm(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Charm", reflect.TypeOf((*MockUniterClient)(nil).Charm), arg0)
}

// CloudAPIVersion mocks base method.
func (m *MockUniterClient) CloudAPIVersion() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudAPIVersion")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudAPIVersion indicates an expected call of CloudAPIVersion.
func (mr *MockUniterClientMockRecorder) CloudAPIVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudAPIVersion", reflect.TypeOf((*MockUniterClient)(nil).CloudAPIVersion))
}

// CloudSpec mocks base method.
func (m *MockUniterClient) CloudSpec() (*params.CloudSpec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudSpec")
	ret0, _ := ret[0].(*params.CloudSpec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudSpec indicates an expected call of CloudSpec.
func (mr *MockUniterClientMockRecorder) CloudSpec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudSpec", reflect.TypeOf((*MockUniterClient)(nil).CloudSpec))
}

// DestroyUnitStorageAttachments mocks base method.
func (m *MockUniterClient) DestroyUnitStorageAttachments(arg0 names.UnitTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyUnitStorageAttachments", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyUnitStorageAttachments indicates an expected call of DestroyUnitStorageAttachments.
func (mr *MockUniterClientMockRecorder) DestroyUnitStorageAttachments(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyUnitStorageAttachments", reflect.TypeOf((*MockUniterClient)(nil).DestroyUnitStorageAttachments), arg0)
}

// GoalState mocks base method.
func (m *MockUniterClient) GoalState() (application.GoalState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GoalState")
	ret0, _ := ret[0].(application.GoalState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GoalState indicates an expected call of GoalState.
func (mr *MockUniterClientMockRecorder) GoalState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoalState", reflect.TypeOf((*MockUniterClient)(nil).GoalState))
}

// LeadershipSettings mocks base method.
func (m *MockUniterClient) LeadershipSettings() uniter.LeadershipSettingsAccessor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeadershipSettings")
	ret0, _ := ret[0].(uniter.LeadershipSettingsAccessor)
	return ret0
}

// LeadershipSettings indicates an expected call of LeadershipSettings.
func (mr *MockUniterClientMockRecorder) LeadershipSettings() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipSettings", reflect.TypeOf((*MockUniterClient)(nil).LeadershipSettings))
}

// Model mocks base method.
func (m *MockUniterClient) Model() (*model.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Model")
	ret0, _ := ret[0].(*model.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Model indicates an expected call of Model.
func (mr *MockUniterClientMockRecorder) Model() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Model", reflect.TypeOf((*MockUniterClient)(nil).Model))
}

// ModelConfig mocks base method.
func (m *MockUniterClient) ModelConfig() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig.
func (mr *MockUniterClientMockRecorder) ModelConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockUniterClient)(nil).ModelConfig))
}

// OpenedMachinePortRangesByEndpoint mocks base method.
func (m *MockUniterClient) OpenedMachinePortRangesByEndpoint(arg0 names.MachineTag) (map[names.UnitTag]network.GroupedPortRanges, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenedMachinePortRangesByEndpoint", arg0)
	ret0, _ := ret[0].(map[names.UnitTag]network.GroupedPortRanges)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenedMachinePortRangesByEndpoint indicates an expected call of OpenedMachinePortRangesByEndpoint.
func (mr *MockUniterClientMockRecorder) OpenedMachinePortRangesByEndpoint(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenedMachinePortRangesByEndpoint", reflect.TypeOf((*MockUniterClient)(nil).OpenedMachinePortRangesByEndpoint), arg0)
}

// OpenedPortRangesByEndpoint mocks base method.
func (m *MockUniterClient) OpenedPortRangesByEndpoint() (map[names.UnitTag]network.GroupedPortRanges, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenedPortRangesByEndpoint")
	ret0, _ := ret[0].(map[names.UnitTag]network.GroupedPortRanges)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenedPortRangesByEndpoint indicates an expected call of OpenedPortRangesByEndpoint.
func (mr *MockUniterClientMockRecorder) OpenedPortRangesByEndpoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenedPortRangesByEndpoint", reflect.TypeOf((*MockUniterClient)(nil).OpenedPortRangesByEndpoint))
}

// Relation mocks base method.
func (m *MockUniterClient) Relation(arg0 names.RelationTag) (Relation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Relation", arg0)
	ret0, _ := ret[0].(Relation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Relation indicates an expected call of Relation.
func (mr *MockUniterClientMockRecorder) Relation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Relation", reflect.TypeOf((*MockUniterClient)(nil).Relation), arg0)
}

// RelationById mocks base method.
func (m *MockUniterClient) RelationById(arg0 int) (Relation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RelationById", arg0)
	ret0, _ := ret[0].(Relation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RelationById indicates an expected call of RelationById.
func (mr *MockUniterClientMockRecorder) RelationById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RelationById", reflect.TypeOf((*MockUniterClient)(nil).RelationById), arg0)
}

// RemoveStorageAttachment mocks base method.
func (m *MockUniterClient) RemoveStorageAttachment(arg0 names.StorageTag, arg1 names.UnitTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveStorageAttachment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveStorageAttachment indicates an expected call of RemoveStorageAttachment.
func (mr *MockUniterClientMockRecorder) RemoveStorageAttachment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveStorageAttachment", reflect.TypeOf((*MockUniterClient)(nil).RemoveStorageAttachment), arg0, arg1)
}

// SLALevel mocks base method.
func (m *MockUniterClient) SLALevel() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SLALevel")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SLALevel indicates an expected call of SLALevel.
func (mr *MockUniterClientMockRecorder) SLALevel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SLALevel", reflect.TypeOf((*MockUniterClient)(nil).SLALevel))
}

// SetUnitWorkloadVersion mocks base method.
func (m *MockUniterClient) SetUnitWorkloadVersion(arg0 names.UnitTag, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUnitWorkloadVersion", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUnitWorkloadVersion indicates an expected call of SetUnitWorkloadVersion.
func (mr *MockUniterClientMockRecorder) SetUnitWorkloadVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUnitWorkloadVersion", reflect.TypeOf((*MockUniterClient)(nil).SetUnitWorkloadVersion), arg0, arg1)
}

// StorageAttachment mocks base method.
func (m *MockUniterClient) StorageAttachment(arg0 names.StorageTag, arg1 names.UnitTag) (params.StorageAttachment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageAttachment", arg0, arg1)
	ret0, _ := ret[0].(params.StorageAttachment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageAttachment indicates an expected call of StorageAttachment.
func (mr *MockUniterClientMockRecorder) StorageAttachment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageAttachment", reflect.TypeOf((*MockUniterClient)(nil).StorageAttachment), arg0, arg1)
}

// StorageAttachmentLife mocks base method.
func (m *MockUniterClient) StorageAttachmentLife(arg0 []params.StorageAttachmentId) ([]params.LifeResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageAttachmentLife", arg0)
	ret0, _ := ret[0].([]params.LifeResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageAttachmentLife indicates an expected call of StorageAttachmentLife.
func (mr *MockUniterClientMockRecorder) StorageAttachmentLife(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageAttachmentLife", reflect.TypeOf((*MockUniterClient)(nil).StorageAttachmentLife), arg0)
}

// Unit mocks base method.
func (m *MockUniterClient) Unit(arg0 names.UnitTag) (Unit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unit", arg0)
	ret0, _ := ret[0].(Unit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unit indicates an expected call of Unit.
func (mr *MockUniterClientMockRecorder) Unit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unit", reflect.TypeOf((*MockUniterClient)(nil).Unit), arg0)
}

// UnitStorageAttachments mocks base method.
func (m *MockUniterClient) UnitStorageAttachments(arg0 names.UnitTag) ([]params.StorageAttachmentId, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnitStorageAttachments", arg0)
	ret0, _ := ret[0].([]params.StorageAttachmentId)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnitStorageAttachments indicates an expected call of UnitStorageAttachments.
func (mr *MockUniterClientMockRecorder) UnitStorageAttachments(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnitStorageAttachments", reflect.TypeOf((*MockUniterClient)(nil).UnitStorageAttachments), arg0)
}

// UnitWorkloadVersion mocks base method.
func (m *MockUniterClient) UnitWorkloadVersion(arg0 names.UnitTag) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnitWorkloadVersion", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnitWorkloadVersion indicates an expected call of UnitWorkloadVersion.
func (mr *MockUniterClientMockRecorder) UnitWorkloadVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnitWorkloadVersion", reflect.TypeOf((*MockUniterClient)(nil).UnitWorkloadVersion), arg0)
}

// UpdateStatusHookInterval mocks base method.
func (m *MockUniterClient) UpdateStatusHookInterval() (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatusHookInterval")
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatusHookInterval indicates an expected call of UpdateStatusHookInterval.
func (mr *MockUniterClientMockRecorder) UpdateStatusHookInterval() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatusHookInterval", reflect.TypeOf((*MockUniterClient)(nil).UpdateStatusHookInterval))
}

// WatchRelationUnits mocks base method.
func (m *MockUniterClient) WatchRelationUnits(arg0 names.RelationTag, arg1 names.UnitTag) (watcher.Watcher[watcher.RelationUnitsChange], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchRelationUnits", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[watcher.RelationUnitsChange])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchRelationUnits indicates an expected call of WatchRelationUnits.
func (mr *MockUniterClientMockRecorder) WatchRelationUnits(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchRelationUnits", reflect.TypeOf((*MockUniterClient)(nil).WatchRelationUnits), arg0, arg1)
}

// WatchStorageAttachment mocks base method.
func (m *MockUniterClient) WatchStorageAttachment(arg0 names.StorageTag, arg1 names.UnitTag) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchStorageAttachment", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchStorageAttachment indicates an expected call of WatchStorageAttachment.
func (mr *MockUniterClientMockRecorder) WatchStorageAttachment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchStorageAttachment", reflect.TypeOf((*MockUniterClient)(nil).WatchStorageAttachment), arg0, arg1)
}

// WatchUpdateStatusHookInterval mocks base method.
func (m *MockUniterClient) WatchUpdateStatusHookInterval() (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchUpdateStatusHookInterval")
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchUpdateStatusHookInterval indicates an expected call of WatchUpdateStatusHookInterval.
func (mr *MockUniterClientMockRecorder) WatchUpdateStatusHookInterval() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchUpdateStatusHookInterval", reflect.TypeOf((*MockUniterClient)(nil).WatchUpdateStatusHookInterval))
}