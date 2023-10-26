// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/common/secrets (interfaces: Model,SecretsConsumer,SecretsMetaState,SecretsRemoveState)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	secrets "github.com/juju/juju/core/secrets"
	config "github.com/juju/juju/environs/config"
	state "github.com/juju/juju/state"
	names "github.com/juju/names/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockModel is a mock of Model interface.
type MockModel struct {
	ctrl     *gomock.Controller
	recorder *MockModelMockRecorder
}

// MockModelMockRecorder is the mock recorder for MockModel.
type MockModelMockRecorder struct {
	mock *MockModel
}

// NewMockModel creates a new mock instance.
func NewMockModel(ctrl *gomock.Controller) *MockModel {
	mock := &MockModel{ctrl: ctrl}
	mock.recorder = &MockModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModel) EXPECT() *MockModelMockRecorder {
	return m.recorder
}

// CloudCredentialTag mocks base method.
func (m *MockModel) CloudCredentialTag() (names.CloudCredentialTag, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredentialTag")
	ret0, _ := ret[0].(names.CloudCredentialTag)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// CloudCredentialTag indicates an expected call of CloudCredentialTag.
func (mr *MockModelMockRecorder) CloudCredentialTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredentialTag", reflect.TypeOf((*MockModel)(nil).CloudCredentialTag))
}

// CloudName mocks base method.
func (m *MockModel) CloudName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudName")
	ret0, _ := ret[0].(string)
	return ret0
}

// CloudName indicates an expected call of CloudName.
func (mr *MockModelMockRecorder) CloudName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudName", reflect.TypeOf((*MockModel)(nil).CloudName))
}

// Config mocks base method.
func (m *MockModel) Config() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Config indicates an expected call of Config.
func (mr *MockModelMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockModel)(nil).Config))
}

// ControllerUUID mocks base method.
func (m *MockModel) ControllerUUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerUUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ControllerUUID indicates an expected call of ControllerUUID.
func (mr *MockModelMockRecorder) ControllerUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerUUID", reflect.TypeOf((*MockModel)(nil).ControllerUUID))
}

// ModelConfig mocks base method.
func (m *MockModel) ModelConfig(arg0 context.Context) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig", arg0)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig.
func (mr *MockModelMockRecorder) ModelConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockModel)(nil).ModelConfig), arg0)
}

// Name mocks base method.
func (m *MockModel) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockModelMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockModel)(nil).Name))
}

// State mocks base method.
func (m *MockModel) State() *state.State {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(*state.State)
	return ret0
}

// State indicates an expected call of State.
func (mr *MockModelMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockModel)(nil).State))
}

// Type mocks base method.
func (m *MockModel) Type() state.ModelType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(state.ModelType)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockModelMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockModel)(nil).Type))
}

// UUID mocks base method.
func (m *MockModel) UUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// UUID indicates an expected call of UUID.
func (mr *MockModelMockRecorder) UUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UUID", reflect.TypeOf((*MockModel)(nil).UUID))
}

// WatchForModelConfigChanges mocks base method.
func (m *MockModel) WatchForModelConfigChanges() state.NotifyWatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchForModelConfigChanges")
	ret0, _ := ret[0].(state.NotifyWatcher)
	return ret0
}

// WatchForModelConfigChanges indicates an expected call of WatchForModelConfigChanges.
func (mr *MockModelMockRecorder) WatchForModelConfigChanges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchForModelConfigChanges", reflect.TypeOf((*MockModel)(nil).WatchForModelConfigChanges))
}

// MockSecretsConsumer is a mock of SecretsConsumer interface.
type MockSecretsConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsConsumerMockRecorder
}

// MockSecretsConsumerMockRecorder is the mock recorder for MockSecretsConsumer.
type MockSecretsConsumerMockRecorder struct {
	mock *MockSecretsConsumer
}

// NewMockSecretsConsumer creates a new mock instance.
func NewMockSecretsConsumer(ctrl *gomock.Controller) *MockSecretsConsumer {
	mock := &MockSecretsConsumer{ctrl: ctrl}
	mock.recorder = &MockSecretsConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretsConsumer) EXPECT() *MockSecretsConsumerMockRecorder {
	return m.recorder
}

// SecretAccess mocks base method.
func (m *MockSecretsConsumer) SecretAccess(arg0 *secrets.URI, arg1 names.Tag) (secrets.SecretRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretAccess", arg0, arg1)
	ret0, _ := ret[0].(secrets.SecretRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SecretAccess indicates an expected call of SecretAccess.
func (mr *MockSecretsConsumerMockRecorder) SecretAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretAccess", reflect.TypeOf((*MockSecretsConsumer)(nil).SecretAccess), arg0, arg1)
}

// MockSecretsMetaState is a mock of SecretsMetaState interface.
type MockSecretsMetaState struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsMetaStateMockRecorder
}

// MockSecretsMetaStateMockRecorder is the mock recorder for MockSecretsMetaState.
type MockSecretsMetaStateMockRecorder struct {
	mock *MockSecretsMetaState
}

// NewMockSecretsMetaState creates a new mock instance.
func NewMockSecretsMetaState(ctrl *gomock.Controller) *MockSecretsMetaState {
	mock := &MockSecretsMetaState{ctrl: ctrl}
	mock.recorder = &MockSecretsMetaStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretsMetaState) EXPECT() *MockSecretsMetaStateMockRecorder {
	return m.recorder
}

// ChangeSecretBackend mocks base method.
func (m *MockSecretsMetaState) ChangeSecretBackend(arg0 state.ChangeSecretBackendParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeSecretBackend", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeSecretBackend indicates an expected call of ChangeSecretBackend.
func (mr *MockSecretsMetaStateMockRecorder) ChangeSecretBackend(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeSecretBackend", reflect.TypeOf((*MockSecretsMetaState)(nil).ChangeSecretBackend), arg0)
}

// ListSecretRevisions mocks base method.
func (m *MockSecretsMetaState) ListSecretRevisions(arg0 *secrets.URI) ([]*secrets.SecretRevisionMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecretRevisions", arg0)
	ret0, _ := ret[0].([]*secrets.SecretRevisionMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecretRevisions indicates an expected call of ListSecretRevisions.
func (mr *MockSecretsMetaStateMockRecorder) ListSecretRevisions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecretRevisions", reflect.TypeOf((*MockSecretsMetaState)(nil).ListSecretRevisions), arg0)
}

// ListSecrets mocks base method.
func (m *MockSecretsMetaState) ListSecrets(arg0 state.SecretsFilter) ([]*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecrets", arg0)
	ret0, _ := ret[0].([]*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecrets indicates an expected call of ListSecrets.
func (mr *MockSecretsMetaStateMockRecorder) ListSecrets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecrets", reflect.TypeOf((*MockSecretsMetaState)(nil).ListSecrets), arg0)
}

// MockSecretsRemoveState is a mock of SecretsRemoveState interface.
type MockSecretsRemoveState struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsRemoveStateMockRecorder
}

// MockSecretsRemoveStateMockRecorder is the mock recorder for MockSecretsRemoveState.
type MockSecretsRemoveStateMockRecorder struct {
	mock *MockSecretsRemoveState
}

// NewMockSecretsRemoveState creates a new mock instance.
func NewMockSecretsRemoveState(ctrl *gomock.Controller) *MockSecretsRemoveState {
	mock := &MockSecretsRemoveState{ctrl: ctrl}
	mock.recorder = &MockSecretsRemoveStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretsRemoveState) EXPECT() *MockSecretsRemoveStateMockRecorder {
	return m.recorder
}

// DeleteSecret mocks base method.
func (m *MockSecretsRemoveState) DeleteSecret(arg0 *secrets.URI, arg1 ...int) ([]secrets.ValueRef, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSecret", varargs...)
	ret0, _ := ret[0].([]secrets.ValueRef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockSecretsRemoveStateMockRecorder) DeleteSecret(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecretsRemoveState)(nil).DeleteSecret), varargs...)
}

// GetSecret mocks base method.
func (m *MockSecretsRemoveState) GetSecret(arg0 *secrets.URI) (*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", arg0)
	ret0, _ := ret[0].(*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockSecretsRemoveStateMockRecorder) GetSecret(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecretsRemoveState)(nil).GetSecret), arg0)
}

// GetSecretRevision mocks base method.
func (m *MockSecretsRemoveState) GetSecretRevision(arg0 *secrets.URI, arg1 int) (*secrets.SecretRevisionMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretRevision", arg0, arg1)
	ret0, _ := ret[0].(*secrets.SecretRevisionMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretRevision indicates an expected call of GetSecretRevision.
func (mr *MockSecretsRemoveStateMockRecorder) GetSecretRevision(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretRevision", reflect.TypeOf((*MockSecretsRemoveState)(nil).GetSecretRevision), arg0, arg1)
}

// ListSecretRevisions mocks base method.
func (m *MockSecretsRemoveState) ListSecretRevisions(arg0 *secrets.URI) ([]*secrets.SecretRevisionMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecretRevisions", arg0)
	ret0, _ := ret[0].([]*secrets.SecretRevisionMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecretRevisions indicates an expected call of ListSecretRevisions.
func (mr *MockSecretsRemoveStateMockRecorder) ListSecretRevisions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecretRevisions", reflect.TypeOf((*MockSecretsRemoveState)(nil).ListSecretRevisions), arg0)
}
