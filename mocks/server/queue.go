// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Queue is an autogenerated mock type for the queue type
type Queue[Element interface{}] struct {
	mock.Mock
}

type Queue_Expecter[Element interface{}] struct {
	mock *mock.Mock
}

func (_m *Queue[Element]) EXPECT() *Queue_Expecter[Element] {
	return &Queue_Expecter[Element]{mock: &_m.Mock}
}

// Append provides a mock function with given fields: element
func (_m *Queue[Element]) Append(element Element) uint64 {
	ret := _m.Called(element)

	if len(ret) == 0 {
		panic("no return value specified for Append")
	}

	var r0 uint64
	if rf, ok := ret.Get(0).(func(Element) uint64); ok {
		r0 = rf(element)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// Queue_Append_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Append'
type Queue_Append_Call[Element interface{}] struct {
	*mock.Call
}

// Append is a helper method to define mock.On call
//   - element Element
func (_e *Queue_Expecter[Element]) Append(element interface{}) *Queue_Append_Call[Element] {
	return &Queue_Append_Call[Element]{Call: _e.mock.On("Append", element)}
}

func (_c *Queue_Append_Call[Element]) Run(run func(element Element)) *Queue_Append_Call[Element] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Element))
	})
	return _c
}

func (_c *Queue_Append_Call[Element]) Return(_a0 uint64) *Queue_Append_Call[Element] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Queue_Append_Call[Element]) RunAndReturn(run func(Element) uint64) *Queue_Append_Call[Element] {
	_c.Call.Return(run)
	return _c
}

// GetLast provides a mock function with given fields:
func (_m *Queue[Element]) GetLast() (Element, uint64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetLast")
	}

	var r0 Element
	var r1 uint64
	var r2 error
	if rf, ok := ret.Get(0).(func() (Element, uint64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() Element); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(Element)
	}

	if rf, ok := ret.Get(1).(func() uint64); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Queue_GetLast_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLast'
type Queue_GetLast_Call[Element interface{}] struct {
	*mock.Call
}

// GetLast is a helper method to define mock.On call
func (_e *Queue_Expecter[Element]) GetLast() *Queue_GetLast_Call[Element] {
	return &Queue_GetLast_Call[Element]{Call: _e.mock.On("GetLast")}
}

func (_c *Queue_GetLast_Call[Element]) Run(run func()) *Queue_GetLast_Call[Element] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Queue_GetLast_Call[Element]) Return(_a0 Element, _a1 uint64, _a2 error) *Queue_GetLast_Call[Element] {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *Queue_GetLast_Call[Element]) RunAndReturn(run func() (Element, uint64, error)) *Queue_GetLast_Call[Element] {
	_c.Call.Return(run)
	return _c
}

// NewQueue creates a new instance of Queue. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQueue[Element interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *Queue[Element] {
	mock := &Queue[Element]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
