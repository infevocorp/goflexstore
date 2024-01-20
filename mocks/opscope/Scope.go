// Code generated by mockery v2.40.1. DO NOT EDIT.

package mockopscope

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Scope is an autogenerated mock type for the Scope type
type Scope struct {
	mock.Mock
}

type Scope_Expecter struct {
	mock *mock.Mock
}

func (_m *Scope) EXPECT() *Scope_Expecter {
	return &Scope_Expecter{mock: &_m.Mock}
}

// Begin provides a mock function with given fields: ctx
func (_m *Scope) Begin(ctx context.Context) (context.Context, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Begin")
	}

	var r0 context.Context
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (context.Context, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) context.Context); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Scope_Begin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Begin'
type Scope_Begin_Call struct {
	*mock.Call
}

// Begin is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Scope_Expecter) Begin(ctx interface{}) *Scope_Begin_Call {
	return &Scope_Begin_Call{Call: _e.mock.On("Begin", ctx)}
}

func (_c *Scope_Begin_Call) Run(run func(ctx context.Context)) *Scope_Begin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Scope_Begin_Call) Return(_a0 context.Context, _a1 error) *Scope_Begin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Scope_Begin_Call) RunAndReturn(run func(context.Context) (context.Context, error)) *Scope_Begin_Call {
	_c.Call.Return(run)
	return _c
}

// End provides a mock function with given fields: ctx, err
func (_m *Scope) End(ctx context.Context, err error) error {
	ret := _m.Called(ctx, err)

	if len(ret) == 0 {
		panic("no return value specified for End")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, error) error); ok {
		r0 = rf(ctx, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Scope_End_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'End'
type Scope_End_Call struct {
	*mock.Call
}

// End is a helper method to define mock.On call
//   - ctx context.Context
//   - err error
func (_e *Scope_Expecter) End(ctx interface{}, err interface{}) *Scope_End_Call {
	return &Scope_End_Call{Call: _e.mock.On("End", ctx, err)}
}

func (_c *Scope_End_Call) Run(run func(ctx context.Context, err error)) *Scope_End_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(error))
	})
	return _c
}

func (_c *Scope_End_Call) Return(_a0 error) *Scope_End_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Scope_End_Call) RunAndReturn(run func(context.Context, error) error) *Scope_End_Call {
	_c.Call.Return(run)
	return _c
}

// EndWithRecover provides a mock function with given fields: ctx, err
func (_m *Scope) EndWithRecover(ctx context.Context, err *error) {
	_m.Called(ctx, err)
}

// Scope_EndWithRecover_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EndWithRecover'
type Scope_EndWithRecover_Call struct {
	*mock.Call
}

// EndWithRecover is a helper method to define mock.On call
//   - ctx context.Context
//   - err *error
func (_e *Scope_Expecter) EndWithRecover(ctx interface{}, err interface{}) *Scope_EndWithRecover_Call {
	return &Scope_EndWithRecover_Call{Call: _e.mock.On("EndWithRecover", ctx, err)}
}

func (_c *Scope_EndWithRecover_Call) Run(run func(ctx context.Context, err *error)) *Scope_EndWithRecover_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*error))
	})
	return _c
}

func (_c *Scope_EndWithRecover_Call) Return() *Scope_EndWithRecover_Call {
	_c.Call.Return()
	return _c
}

func (_c *Scope_EndWithRecover_Call) RunAndReturn(run func(context.Context, *error)) *Scope_EndWithRecover_Call {
	_c.Call.Return(run)
	return _c
}

// NewScope creates a new instance of Scope. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScope(t interface {
	mock.TestingT
	Cleanup(func())
}) *Scope {
	mock := &Scope{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}