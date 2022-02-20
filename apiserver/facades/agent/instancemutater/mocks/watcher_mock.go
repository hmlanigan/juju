// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/state (interfaces: NotifyWatcher,StringsWatcher)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockNotifyWatcher is a mock of NotifyWatcher interface.
type MockNotifyWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockNotifyWatcherMockRecorder
}

// MockNotifyWatcherMockRecorder is the mock recorder for MockNotifyWatcher.
type MockNotifyWatcherMockRecorder struct {
	mock *MockNotifyWatcher
}

// NewMockNotifyWatcher creates a new mock instance.
func NewMockNotifyWatcher(ctrl *gomock.Controller) *MockNotifyWatcher {
	mock := &MockNotifyWatcher{ctrl: ctrl}
	mock.recorder = &MockNotifyWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotifyWatcher) EXPECT() *MockNotifyWatcherMockRecorder {
	return m.recorder
}

// Changes mocks base method.
func (m *MockNotifyWatcher) Changes() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Changes")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Changes indicates an expected call of Changes.
func (mr *MockNotifyWatcherMockRecorder) Changes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Changes", reflect.TypeOf((*MockNotifyWatcher)(nil).Changes))
}

// Err mocks base method.
func (m *MockNotifyWatcher) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockNotifyWatcherMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockNotifyWatcher)(nil).Err))
}

// Kill mocks base method.
func (m *MockNotifyWatcher) Kill() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill.
func (mr *MockNotifyWatcherMockRecorder) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockNotifyWatcher)(nil).Kill))
}

// Stop mocks base method.
func (m *MockNotifyWatcher) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockNotifyWatcherMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockNotifyWatcher)(nil).Stop))
}

// Wait mocks base method.
func (m *MockNotifyWatcher) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockNotifyWatcherMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockNotifyWatcher)(nil).Wait))
}

// MockStringsWatcher is a mock of StringsWatcher interface.
type MockStringsWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockStringsWatcherMockRecorder
}

// MockStringsWatcherMockRecorder is the mock recorder for MockStringsWatcher.
type MockStringsWatcherMockRecorder struct {
	mock *MockStringsWatcher
}

// NewMockStringsWatcher creates a new mock instance.
func NewMockStringsWatcher(ctrl *gomock.Controller) *MockStringsWatcher {
	mock := &MockStringsWatcher{ctrl: ctrl}
	mock.recorder = &MockStringsWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStringsWatcher) EXPECT() *MockStringsWatcherMockRecorder {
	return m.recorder
}

// Changes mocks base method.
func (m *MockStringsWatcher) Changes() <-chan []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Changes")
	ret0, _ := ret[0].(<-chan []string)
	return ret0
}

// Changes indicates an expected call of Changes.
func (mr *MockStringsWatcherMockRecorder) Changes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Changes", reflect.TypeOf((*MockStringsWatcher)(nil).Changes))
}

// Err mocks base method.
func (m *MockStringsWatcher) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockStringsWatcherMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockStringsWatcher)(nil).Err))
}

// Kill mocks base method.
func (m *MockStringsWatcher) Kill() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill.
func (mr *MockStringsWatcherMockRecorder) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockStringsWatcher)(nil).Kill))
}

// Stop mocks base method.
func (m *MockStringsWatcher) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockStringsWatcherMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockStringsWatcher)(nil).Stop))
}

// Wait mocks base method.
func (m *MockStringsWatcher) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockStringsWatcherMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockStringsWatcher)(nil).Wait))
}
