// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/changestream (interfaces: DBGetter,Logger,EventMultiplexerWorker,FileNotifyWatcher)

// Package changestream is a generated GoMock package.
package changestream

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	changestream "github.com/juju/juju/core/changestream"
	database "github.com/juju/juju/core/database"
)

// MockDBGetter is a mock of DBGetter interface.
type MockDBGetter struct {
	ctrl     *gomock.Controller
	recorder *MockDBGetterMockRecorder
}

// MockDBGetterMockRecorder is the mock recorder for MockDBGetter.
type MockDBGetterMockRecorder struct {
	mock *MockDBGetter
}

// NewMockDBGetter creates a new mock instance.
func NewMockDBGetter(ctrl *gomock.Controller) *MockDBGetter {
	mock := &MockDBGetter{ctrl: ctrl}
	mock.recorder = &MockDBGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBGetter) EXPECT() *MockDBGetterMockRecorder {
	return m.recorder
}

// GetDB mocks base method.
func (m *MockDBGetter) GetDB(arg0 string) (database.TxnRunner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDB", arg0)
	ret0, _ := ret[0].(database.TxnRunner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDB indicates an expected call of GetDB.
func (mr *MockDBGetterMockRecorder) GetDB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDB", reflect.TypeOf((*MockDBGetter)(nil).GetDB), arg0)
}

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debugf mocks base method.
func (m *MockLogger) Debugf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf.
func (mr *MockLoggerMockRecorder) Debugf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*MockLogger)(nil).Debugf), varargs...)
}

// Errorf mocks base method.
func (m *MockLogger) Errorf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockLoggerMockRecorder) Errorf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockLogger)(nil).Errorf), varargs...)
}

// Infof mocks base method.
func (m *MockLogger) Infof(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockLoggerMockRecorder) Infof(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockLogger)(nil).Infof), varargs...)
}

// IsTraceEnabled mocks base method.
func (m *MockLogger) IsTraceEnabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTraceEnabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTraceEnabled indicates an expected call of IsTraceEnabled.
func (mr *MockLoggerMockRecorder) IsTraceEnabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTraceEnabled", reflect.TypeOf((*MockLogger)(nil).IsTraceEnabled))
}

// Tracef mocks base method.
func (m *MockLogger) Tracef(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Tracef", varargs...)
}

// Tracef indicates an expected call of Tracef.
func (mr *MockLoggerMockRecorder) Tracef(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tracef", reflect.TypeOf((*MockLogger)(nil).Tracef), varargs...)
}

// Warningf mocks base method.
func (m *MockLogger) Warningf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warningf", varargs...)
}

// Warningf indicates an expected call of Warningf.
func (mr *MockLoggerMockRecorder) Warningf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warningf", reflect.TypeOf((*MockLogger)(nil).Warningf), varargs...)
}

// MockEventMultiplexerWorker is a mock of EventMultiplexerWorker interface.
type MockEventMultiplexerWorker struct {
	ctrl     *gomock.Controller
	recorder *MockEventMultiplexerWorkerMockRecorder
}

// MockEventMultiplexerWorkerMockRecorder is the mock recorder for MockEventMultiplexerWorker.
type MockEventMultiplexerWorkerMockRecorder struct {
	mock *MockEventMultiplexerWorker
}

// NewMockEventMultiplexerWorker creates a new mock instance.
func NewMockEventMultiplexerWorker(ctrl *gomock.Controller) *MockEventMultiplexerWorker {
	mock := &MockEventMultiplexerWorker{ctrl: ctrl}
	mock.recorder = &MockEventMultiplexerWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventMultiplexerWorker) EXPECT() *MockEventMultiplexerWorkerMockRecorder {
	return m.recorder
}

// EventSource mocks base method.
func (m *MockEventMultiplexerWorker) EventSource() changestream.EventSource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventSource")
	ret0, _ := ret[0].(changestream.EventSource)
	return ret0
}

// EventSource indicates an expected call of EventSource.
func (mr *MockEventMultiplexerWorkerMockRecorder) EventSource() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventSource", reflect.TypeOf((*MockEventMultiplexerWorker)(nil).EventSource))
}

// Kill mocks base method.
func (m *MockEventMultiplexerWorker) Kill() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill.
func (mr *MockEventMultiplexerWorkerMockRecorder) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockEventMultiplexerWorker)(nil).Kill))
}

// Wait mocks base method.
func (m *MockEventMultiplexerWorker) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockEventMultiplexerWorkerMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockEventMultiplexerWorker)(nil).Wait))
}

// MockFileNotifyWatcher is a mock of FileNotifyWatcher interface.
type MockFileNotifyWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockFileNotifyWatcherMockRecorder
}

// MockFileNotifyWatcherMockRecorder is the mock recorder for MockFileNotifyWatcher.
type MockFileNotifyWatcherMockRecorder struct {
	mock *MockFileNotifyWatcher
}

// NewMockFileNotifyWatcher creates a new mock instance.
func NewMockFileNotifyWatcher(ctrl *gomock.Controller) *MockFileNotifyWatcher {
	mock := &MockFileNotifyWatcher{ctrl: ctrl}
	mock.recorder = &MockFileNotifyWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileNotifyWatcher) EXPECT() *MockFileNotifyWatcherMockRecorder {
	return m.recorder
}

// Changes mocks base method.
func (m *MockFileNotifyWatcher) Changes(arg0 string) (<-chan bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Changes", arg0)
	ret0, _ := ret[0].(<-chan bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Changes indicates an expected call of Changes.
func (mr *MockFileNotifyWatcherMockRecorder) Changes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Changes", reflect.TypeOf((*MockFileNotifyWatcher)(nil).Changes), arg0)
}