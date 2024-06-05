// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// HttpHandler is an autogenerated mock type for the httpHandler type
type HttpHandler struct {
	mock.Mock
}

type HttpHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *HttpHandler) EXPECT() *HttpHandler_Expecter {
	return &HttpHandler_Expecter{mock: &_m.Mock}
}

// ServeHTTP provides a mock function with given fields: _a0, _a1
func (_m *HttpHandler) ServeHTTP(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// HttpHandler_ServeHTTP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServeHTTP'
type HttpHandler_ServeHTTP_Call struct {
	*mock.Call
}

// ServeHTTP is a helper method to define mock.On call
//   - _a0 http.ResponseWriter
//   - _a1 *http.Request
func (_e *HttpHandler_Expecter) ServeHTTP(_a0 interface{}, _a1 interface{}) *HttpHandler_ServeHTTP_Call {
	return &HttpHandler_ServeHTTP_Call{Call: _e.mock.On("ServeHTTP", _a0, _a1)}
}

func (_c *HttpHandler_ServeHTTP_Call) Run(run func(_a0 http.ResponseWriter, _a1 *http.Request)) *HttpHandler_ServeHTTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *HttpHandler_ServeHTTP_Call) Return() *HttpHandler_ServeHTTP_Call {
	_c.Call.Return()
	return _c
}

func (_c *HttpHandler_ServeHTTP_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *HttpHandler_ServeHTTP_Call {
	_c.Call.Return(run)
	return _c
}

// NewHttpHandler creates a new instance of HttpHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHttpHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *HttpHandler {
	mock := &HttpHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
