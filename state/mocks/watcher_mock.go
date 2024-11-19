// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/state/watcher (interfaces: BaseWatcher)
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/watcher_mock.go github.com/juju/juju/state/watcher BaseWatcher
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	watcher "github.com/juju/juju/state/watcher"
	gomock "go.uber.org/mock/gomock"
)

// MockBaseWatcher is a mock of BaseWatcher interface.
type MockBaseWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockBaseWatcherMockRecorder
}

// MockBaseWatcherMockRecorder is the mock recorder for MockBaseWatcher.
type MockBaseWatcherMockRecorder struct {
	mock *MockBaseWatcher
}

// NewMockBaseWatcher creates a new mock instance.
func NewMockBaseWatcher(ctrl *gomock.Controller) *MockBaseWatcher {
	mock := &MockBaseWatcher{ctrl: ctrl}
	mock.recorder = &MockBaseWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBaseWatcher) EXPECT() *MockBaseWatcherMockRecorder {
	return m.recorder
}

// Dead mocks base method.
func (m *MockBaseWatcher) Dead() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dead")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Dead indicates an expected call of Dead.
func (mr *MockBaseWatcherMockRecorder) Dead() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dead", reflect.TypeOf((*MockBaseWatcher)(nil).Dead))
}

// Err mocks base method.
func (m *MockBaseWatcher) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockBaseWatcherMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockBaseWatcher)(nil).Err))
}

// Kill mocks base method.
func (m *MockBaseWatcher) Kill() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill.
func (mr *MockBaseWatcherMockRecorder) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockBaseWatcher)(nil).Kill))
}

// Unwatch mocks base method.
func (m *MockBaseWatcher) Unwatch(arg0 string, arg1 any, arg2 chan<- watcher.Change) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unwatch", arg0, arg1, arg2)
}

// Unwatch indicates an expected call of Unwatch.
func (mr *MockBaseWatcherMockRecorder) Unwatch(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unwatch", reflect.TypeOf((*MockBaseWatcher)(nil).Unwatch), arg0, arg1, arg2)
}

// UnwatchCollection mocks base method.
func (m *MockBaseWatcher) UnwatchCollection(arg0 string, arg1 chan<- watcher.Change) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnwatchCollection", arg0, arg1)
}

// UnwatchCollection indicates an expected call of UnwatchCollection.
func (mr *MockBaseWatcherMockRecorder) UnwatchCollection(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnwatchCollection", reflect.TypeOf((*MockBaseWatcher)(nil).UnwatchCollection), arg0, arg1)
}

// Wait mocks base method.
func (m *MockBaseWatcher) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockBaseWatcherMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockBaseWatcher)(nil).Wait))
}

// Watch mocks base method.
func (m *MockBaseWatcher) Watch(arg0 string, arg1 any, arg2 chan<- watcher.Change) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Watch", arg0, arg1, arg2)
}

// Watch indicates an expected call of Watch.
func (mr *MockBaseWatcherMockRecorder) Watch(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockBaseWatcher)(nil).Watch), arg0, arg1, arg2)
}

// WatchCollection mocks base method.
func (m *MockBaseWatcher) WatchCollection(arg0 string, arg1 chan<- watcher.Change) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WatchCollection", arg0, arg1)
}

// WatchCollection indicates an expected call of WatchCollection.
func (mr *MockBaseWatcherMockRecorder) WatchCollection(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchCollection", reflect.TypeOf((*MockBaseWatcher)(nil).WatchCollection), arg0, arg1)
}

// WatchCollectionWithFilter mocks base method.
func (m *MockBaseWatcher) WatchCollectionWithFilter(arg0 string, arg1 chan<- watcher.Change, arg2 func(any) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WatchCollectionWithFilter", arg0, arg1, arg2)
}

// WatchCollectionWithFilter indicates an expected call of WatchCollectionWithFilter.
func (mr *MockBaseWatcherMockRecorder) WatchCollectionWithFilter(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchCollectionWithFilter", reflect.TypeOf((*MockBaseWatcher)(nil).WatchCollectionWithFilter), arg0, arg1, arg2)
}

// WatchMulti mocks base method.
func (m *MockBaseWatcher) WatchMulti(arg0 string, arg1 []any, arg2 chan<- watcher.Change) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchMulti", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WatchMulti indicates an expected call of WatchMulti.
func (mr *MockBaseWatcherMockRecorder) WatchMulti(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchMulti", reflect.TypeOf((*MockBaseWatcher)(nil).WatchMulti), arg0, arg1, arg2)
}