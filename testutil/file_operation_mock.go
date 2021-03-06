// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ymgyt/cloudops/usecase (interfaces: FileOps)

// Package testutil is a generated GoMock package.
package testutil

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	usecase "github.com/ymgyt/cloudops/usecase"
)

// MockFileOps is a mock of FileOps interface
type MockFileOps struct {
	ctrl     *gomock.Controller
	recorder *MockFileOpsMockRecorder
}

// MockFileOpsMockRecorder is the mock recorder for MockFileOps
type MockFileOpsMockRecorder struct {
	mock *MockFileOps
}

// NewMockFileOps creates a new mock instance
func NewMockFileOps(ctrl *gomock.Controller) *MockFileOps {
	mock := &MockFileOps{ctrl: ctrl}
	mock.recorder = &MockFileOpsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileOps) EXPECT() *MockFileOpsMockRecorder {
	return m.recorder
}

// CopyLocalToRemote mocks base method
func (m *MockFileOps) CopyLocalToRemote(arg0 *usecase.CopyLocalToRemoteInput) (*usecase.CopyLocalToRemoteOutput, error) {
	ret := m.ctrl.Call(m, "CopyLocalToRemote", arg0)
	ret0, _ := ret[0].(*usecase.CopyLocalToRemoteOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CopyLocalToRemote indicates an expected call of CopyLocalToRemote
func (mr *MockFileOpsMockRecorder) CopyLocalToRemote(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyLocalToRemote", reflect.TypeOf((*MockFileOps)(nil).CopyLocalToRemote), arg0)
}
