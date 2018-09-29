// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ymgyt/cloudops/core (interfaces: Confirmer)

// Package testutil is a generated GoMock package.
package testutil

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	core "github.com/ymgyt/cloudops/core"
)

// MockConfirmer is a mock of Confirmer interface
type MockConfirmer struct {
	ctrl     *gomock.Controller
	recorder *MockConfirmerMockRecorder
}

// MockConfirmerMockRecorder is the mock recorder for MockConfirmer
type MockConfirmerMockRecorder struct {
	mock *MockConfirmer
}

// NewMockConfirmer creates a new mock instance
func NewMockConfirmer(ctrl *gomock.Controller) *MockConfirmer {
	mock := &MockConfirmer{ctrl: ctrl}
	mock.recorder = &MockConfirmerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfirmer) EXPECT() *MockConfirmerMockRecorder {
	return m.recorder
}

// Confirm mocks base method
func (m *MockConfirmer) Confirm(arg0 string, arg1 core.Resources) (bool, error) {
	ret := m.ctrl.Call(m, "Confirm", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Confirm indicates an expected call of Confirm
func (mr *MockConfirmerMockRecorder) Confirm(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Confirm", reflect.TypeOf((*MockConfirmer)(nil).Confirm), arg0, arg1)
}
