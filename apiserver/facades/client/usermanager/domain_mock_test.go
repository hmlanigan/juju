// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/usermanager (interfaces: AccessService)
//
// Generated by this command:
//
//	mockgen -typed -package usermanager_test -destination domain_mock_test.go github.com/juju/juju/apiserver/facades/client/usermanager AccessService
//

// Package usermanager_test is a generated GoMock package.
package usermanager_test

import (
	context "context"
	reflect "reflect"

	model "github.com/juju/juju/core/model"
	permission "github.com/juju/juju/core/permission"
	user "github.com/juju/juju/core/user"
	access "github.com/juju/juju/domain/access"
	service "github.com/juju/juju/domain/access/service"
	auth "github.com/juju/juju/internal/auth"
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

// AddUser mocks base method.
func (m *MockAccessService) AddUser(arg0 context.Context, arg1 service.AddUserArg) (user.UUID, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1)
	ret0, _ := ret[0].(user.UUID)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddUser indicates an expected call of AddUser.
func (mr *MockAccessServiceMockRecorder) AddUser(arg0, arg1 any) *MockAccessServiceAddUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockAccessService)(nil).AddUser), arg0, arg1)
	return &MockAccessServiceAddUserCall{Call: call}
}

// MockAccessServiceAddUserCall wrap *gomock.Call
type MockAccessServiceAddUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceAddUserCall) Return(arg0 user.UUID, arg1 []byte, arg2 error) *MockAccessServiceAddUserCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceAddUserCall) Do(f func(context.Context, service.AddUserArg) (user.UUID, []byte, error)) *MockAccessServiceAddUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceAddUserCall) DoAndReturn(f func(context.Context, service.AddUserArg) (user.UUID, []byte, error)) *MockAccessServiceAddUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DisableUserAuthentication mocks base method.
func (m *MockAccessService) DisableUserAuthentication(arg0 context.Context, arg1 user.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableUserAuthentication", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DisableUserAuthentication indicates an expected call of DisableUserAuthentication.
func (mr *MockAccessServiceMockRecorder) DisableUserAuthentication(arg0, arg1 any) *MockAccessServiceDisableUserAuthenticationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableUserAuthentication", reflect.TypeOf((*MockAccessService)(nil).DisableUserAuthentication), arg0, arg1)
	return &MockAccessServiceDisableUserAuthenticationCall{Call: call}
}

// MockAccessServiceDisableUserAuthenticationCall wrap *gomock.Call
type MockAccessServiceDisableUserAuthenticationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceDisableUserAuthenticationCall) Return(arg0 error) *MockAccessServiceDisableUserAuthenticationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceDisableUserAuthenticationCall) Do(f func(context.Context, user.Name) error) *MockAccessServiceDisableUserAuthenticationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceDisableUserAuthenticationCall) DoAndReturn(f func(context.Context, user.Name) error) *MockAccessServiceDisableUserAuthenticationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// EnableUserAuthentication mocks base method.
func (m *MockAccessService) EnableUserAuthentication(arg0 context.Context, arg1 user.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableUserAuthentication", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnableUserAuthentication indicates an expected call of EnableUserAuthentication.
func (mr *MockAccessServiceMockRecorder) EnableUserAuthentication(arg0, arg1 any) *MockAccessServiceEnableUserAuthenticationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableUserAuthentication", reflect.TypeOf((*MockAccessService)(nil).EnableUserAuthentication), arg0, arg1)
	return &MockAccessServiceEnableUserAuthenticationCall{Call: call}
}

// MockAccessServiceEnableUserAuthenticationCall wrap *gomock.Call
type MockAccessServiceEnableUserAuthenticationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceEnableUserAuthenticationCall) Return(arg0 error) *MockAccessServiceEnableUserAuthenticationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceEnableUserAuthenticationCall) Do(f func(context.Context, user.Name) error) *MockAccessServiceEnableUserAuthenticationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceEnableUserAuthenticationCall) DoAndReturn(f func(context.Context, user.Name) error) *MockAccessServiceEnableUserAuthenticationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetAllUsers mocks base method.
func (m *MockAccessService) GetAllUsers(arg0 context.Context, arg1 bool) ([]user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", arg0, arg1)
	ret0, _ := ret[0].([]user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockAccessServiceMockRecorder) GetAllUsers(arg0, arg1 any) *MockAccessServiceGetAllUsersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockAccessService)(nil).GetAllUsers), arg0, arg1)
	return &MockAccessServiceGetAllUsersCall{Call: call}
}

// MockAccessServiceGetAllUsersCall wrap *gomock.Call
type MockAccessServiceGetAllUsersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceGetAllUsersCall) Return(arg0 []user.User, arg1 error) *MockAccessServiceGetAllUsersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceGetAllUsersCall) Do(f func(context.Context, bool) ([]user.User, error)) *MockAccessServiceGetAllUsersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceGetAllUsersCall) DoAndReturn(f func(context.Context, bool) ([]user.User, error)) *MockAccessServiceGetAllUsersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetModelUsers mocks base method.
func (m *MockAccessService) GetModelUsers(arg0 context.Context, arg1 user.Name, arg2 model.UUID) ([]access.ModelUserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModelUsers", arg0, arg1, arg2)
	ret0, _ := ret[0].([]access.ModelUserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetModelUsers indicates an expected call of GetModelUsers.
func (mr *MockAccessServiceMockRecorder) GetModelUsers(arg0, arg1, arg2 any) *MockAccessServiceGetModelUsersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModelUsers", reflect.TypeOf((*MockAccessService)(nil).GetModelUsers), arg0, arg1, arg2)
	return &MockAccessServiceGetModelUsersCall{Call: call}
}

// MockAccessServiceGetModelUsersCall wrap *gomock.Call
type MockAccessServiceGetModelUsersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceGetModelUsersCall) Return(arg0 []access.ModelUserInfo, arg1 error) *MockAccessServiceGetModelUsersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceGetModelUsersCall) Do(f func(context.Context, user.Name, model.UUID) ([]access.ModelUserInfo, error)) *MockAccessServiceGetModelUsersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceGetModelUsersCall) DoAndReturn(f func(context.Context, user.Name, model.UUID) ([]access.ModelUserInfo, error)) *MockAccessServiceGetModelUsersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetUserByName mocks base method.
func (m *MockAccessService) GetUserByName(arg0 context.Context, arg1 user.Name) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByName", arg0, arg1)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByName indicates an expected call of GetUserByName.
func (mr *MockAccessServiceMockRecorder) GetUserByName(arg0, arg1 any) *MockAccessServiceGetUserByNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByName", reflect.TypeOf((*MockAccessService)(nil).GetUserByName), arg0, arg1)
	return &MockAccessServiceGetUserByNameCall{Call: call}
}

// MockAccessServiceGetUserByNameCall wrap *gomock.Call
type MockAccessServiceGetUserByNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceGetUserByNameCall) Return(arg0 user.User, arg1 error) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceGetUserByNameCall) Do(f func(context.Context, user.Name) (user.User, error)) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceGetUserByNameCall) DoAndReturn(f func(context.Context, user.Name) (user.User, error)) *MockAccessServiceGetUserByNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ReadUserAccessLevelForTarget mocks base method.
func (m *MockAccessService) ReadUserAccessLevelForTarget(arg0 context.Context, arg1 user.Name, arg2 permission.ID) (permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserAccessLevelForTarget", arg0, arg1, arg2)
	ret0, _ := ret[0].(permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserAccessLevelForTarget indicates an expected call of ReadUserAccessLevelForTarget.
func (mr *MockAccessServiceMockRecorder) ReadUserAccessLevelForTarget(arg0, arg1, arg2 any) *MockAccessServiceReadUserAccessLevelForTargetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserAccessLevelForTarget", reflect.TypeOf((*MockAccessService)(nil).ReadUserAccessLevelForTarget), arg0, arg1, arg2)
	return &MockAccessServiceReadUserAccessLevelForTargetCall{Call: call}
}

// MockAccessServiceReadUserAccessLevelForTargetCall wrap *gomock.Call
type MockAccessServiceReadUserAccessLevelForTargetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceReadUserAccessLevelForTargetCall) Return(arg0 permission.Access, arg1 error) *MockAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceReadUserAccessLevelForTargetCall) Do(f func(context.Context, user.Name, permission.ID) (permission.Access, error)) *MockAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceReadUserAccessLevelForTargetCall) DoAndReturn(f func(context.Context, user.Name, permission.ID) (permission.Access, error)) *MockAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RemoveUser mocks base method.
func (m *MockAccessService) RemoveUser(arg0 context.Context, arg1 user.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUser indicates an expected call of RemoveUser.
func (mr *MockAccessServiceMockRecorder) RemoveUser(arg0, arg1 any) *MockAccessServiceRemoveUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUser", reflect.TypeOf((*MockAccessService)(nil).RemoveUser), arg0, arg1)
	return &MockAccessServiceRemoveUserCall{Call: call}
}

// MockAccessServiceRemoveUserCall wrap *gomock.Call
type MockAccessServiceRemoveUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceRemoveUserCall) Return(arg0 error) *MockAccessServiceRemoveUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceRemoveUserCall) Do(f func(context.Context, user.Name) error) *MockAccessServiceRemoveUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceRemoveUserCall) DoAndReturn(f func(context.Context, user.Name) error) *MockAccessServiceRemoveUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ResetPassword mocks base method.
func (m *MockAccessService) ResetPassword(arg0 context.Context, arg1 user.Name) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockAccessServiceMockRecorder) ResetPassword(arg0, arg1 any) *MockAccessServiceResetPasswordCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockAccessService)(nil).ResetPassword), arg0, arg1)
	return &MockAccessServiceResetPasswordCall{Call: call}
}

// MockAccessServiceResetPasswordCall wrap *gomock.Call
type MockAccessServiceResetPasswordCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceResetPasswordCall) Return(arg0 []byte, arg1 error) *MockAccessServiceResetPasswordCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceResetPasswordCall) Do(f func(context.Context, user.Name) ([]byte, error)) *MockAccessServiceResetPasswordCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceResetPasswordCall) DoAndReturn(f func(context.Context, user.Name) ([]byte, error)) *MockAccessServiceResetPasswordCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetPassword mocks base method.
func (m *MockAccessService) SetPassword(arg0 context.Context, arg1 user.Name, arg2 auth.Password) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPassword indicates an expected call of SetPassword.
func (mr *MockAccessServiceMockRecorder) SetPassword(arg0, arg1, arg2 any) *MockAccessServiceSetPasswordCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPassword", reflect.TypeOf((*MockAccessService)(nil).SetPassword), arg0, arg1, arg2)
	return &MockAccessServiceSetPasswordCall{Call: call}
}

// MockAccessServiceSetPasswordCall wrap *gomock.Call
type MockAccessServiceSetPasswordCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockAccessServiceSetPasswordCall) Return(arg0 error) *MockAccessServiceSetPasswordCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockAccessServiceSetPasswordCall) Do(f func(context.Context, user.Name, auth.Password) error) *MockAccessServiceSetPasswordCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockAccessServiceSetPasswordCall) DoAndReturn(f func(context.Context, user.Name, auth.Password) error) *MockAccessServiceSetPasswordCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
