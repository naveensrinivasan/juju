// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/api/base (interfaces: APICallCloser,ClientFacade)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	httprequest "github.com/juju/httprequest"
	base "github.com/juju/juju/api/base"
	names_v3 "gopkg.in/juju/names.v3"
	httpbakery "gopkg.in/macaroon-bakery.v2-unstable/httpbakery"
	http "net/http"
	url "net/url"
	reflect "reflect"
)

// MockAPICallCloser is a mock of APICallCloser interface
type MockAPICallCloser struct {
	ctrl     *gomock.Controller
	recorder *MockAPICallCloserMockRecorder
}

// MockAPICallCloserMockRecorder is the mock recorder for MockAPICallCloser
type MockAPICallCloserMockRecorder struct {
	mock *MockAPICallCloser
}

// NewMockAPICallCloser creates a new mock instance
func NewMockAPICallCloser(ctrl *gomock.Controller) *MockAPICallCloser {
	mock := &MockAPICallCloser{ctrl: ctrl}
	mock.recorder = &MockAPICallCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAPICallCloser) EXPECT() *MockAPICallCloserMockRecorder {
	return m.recorder
}

// APICall mocks base method
func (m *MockAPICallCloser) APICall(arg0 string, arg1 int, arg2, arg3 string, arg4, arg5 interface{}) error {
	ret := m.ctrl.Call(m, "APICall", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// APICall indicates an expected call of APICall
func (mr *MockAPICallCloserMockRecorder) APICall(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APICall", reflect.TypeOf((*MockAPICallCloser)(nil).APICall), arg0, arg1, arg2, arg3, arg4, arg5)
}

// BakeryClient mocks base method
func (m *MockAPICallCloser) BakeryClient() *httpbakery.Client {
	ret := m.ctrl.Call(m, "BakeryClient")
	ret0, _ := ret[0].(*httpbakery.Client)
	return ret0
}

// BakeryClient indicates an expected call of BakeryClient
func (mr *MockAPICallCloserMockRecorder) BakeryClient() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BakeryClient", reflect.TypeOf((*MockAPICallCloser)(nil).BakeryClient))
}

// BestFacadeVersion mocks base method
func (m *MockAPICallCloser) BestFacadeVersion(arg0 string) int {
	ret := m.ctrl.Call(m, "BestFacadeVersion", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// BestFacadeVersion indicates an expected call of BestFacadeVersion
func (mr *MockAPICallCloserMockRecorder) BestFacadeVersion(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestFacadeVersion", reflect.TypeOf((*MockAPICallCloser)(nil).BestFacadeVersion), arg0)
}

// Close mocks base method
func (m *MockAPICallCloser) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockAPICallCloserMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAPICallCloser)(nil).Close))
}

// ConnectControllerStream mocks base method
func (m *MockAPICallCloser) ConnectControllerStream(arg0 string, arg1 url.Values, arg2 http.Header) (base.Stream, error) {
	ret := m.ctrl.Call(m, "ConnectControllerStream", arg0, arg1, arg2)
	ret0, _ := ret[0].(base.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectControllerStream indicates an expected call of ConnectControllerStream
func (mr *MockAPICallCloserMockRecorder) ConnectControllerStream(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectControllerStream", reflect.TypeOf((*MockAPICallCloser)(nil).ConnectControllerStream), arg0, arg1, arg2)
}

// ConnectStream mocks base method
func (m *MockAPICallCloser) ConnectStream(arg0 string, arg1 url.Values) (base.Stream, error) {
	ret := m.ctrl.Call(m, "ConnectStream", arg0, arg1)
	ret0, _ := ret[0].(base.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectStream indicates an expected call of ConnectStream
func (mr *MockAPICallCloserMockRecorder) ConnectStream(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectStream", reflect.TypeOf((*MockAPICallCloser)(nil).ConnectStream), arg0, arg1)
}

// HTTPClient mocks base method
func (m *MockAPICallCloser) HTTPClient() (*httprequest.Client, error) {
	ret := m.ctrl.Call(m, "HTTPClient")
	ret0, _ := ret[0].(*httprequest.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HTTPClient indicates an expected call of HTTPClient
func (mr *MockAPICallCloserMockRecorder) HTTPClient() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPClient", reflect.TypeOf((*MockAPICallCloser)(nil).HTTPClient))
}

// ModelTag mocks base method
func (m *MockAPICallCloser) ModelTag() (names_v3.ModelTag, bool) {
	ret := m.ctrl.Call(m, "ModelTag")
	ret0, _ := ret[0].(names_v3.ModelTag)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ModelTag indicates an expected call of ModelTag
func (mr *MockAPICallCloserMockRecorder) ModelTag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelTag", reflect.TypeOf((*MockAPICallCloser)(nil).ModelTag))
}

// MockClientFacade is a mock of ClientFacade interface
type MockClientFacade struct {
	ctrl     *gomock.Controller
	recorder *MockClientFacadeMockRecorder
}

// MockClientFacadeMockRecorder is the mock recorder for MockClientFacade
type MockClientFacadeMockRecorder struct {
	mock *MockClientFacade
}

// NewMockClientFacade creates a new mock instance
func NewMockClientFacade(ctrl *gomock.Controller) *MockClientFacade {
	mock := &MockClientFacade{ctrl: ctrl}
	mock.recorder = &MockClientFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClientFacade) EXPECT() *MockClientFacadeMockRecorder {
	return m.recorder
}

// BestAPIVersion mocks base method
func (m *MockClientFacade) BestAPIVersion() int {
	ret := m.ctrl.Call(m, "BestAPIVersion")
	ret0, _ := ret[0].(int)
	return ret0
}

// BestAPIVersion indicates an expected call of BestAPIVersion
func (mr *MockClientFacadeMockRecorder) BestAPIVersion() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestAPIVersion", reflect.TypeOf((*MockClientFacade)(nil).BestAPIVersion))
}

// Close mocks base method
func (m *MockClientFacade) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockClientFacadeMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClientFacade)(nil).Close))
}
