// Code generated by mockery v2.46.0. DO NOT EDIT.

package mockusecase

import (
	context "context"
	dto "date-apps-be/internal/usecase/user/dto"

	mock "github.com/stretchr/testify/mock"

	model "date-apps-be/internal/model"
)

// UserUsecase is an autogenerated mock type for the UserUsecase type
type UserUsecase struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *UserUsecase) CreateUser(ctx context.Context, user *dto.CreateUser) (string, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.CreateUser) (string, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.CreateUser) string); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.CreateUser) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, userUID
func (_m *UserUsecase) GetUser(ctx context.Context, userUID string) (*model.User, error) {
	ret := _m.Called(ctx, userUID)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, userUID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, userUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmailOrPhoneNumber provides a mock function with given fields: ctx, email, phoneNumber
func (_m *UserUsecase) GetUserByEmailOrPhoneNumber(ctx context.Context, email string, phoneNumber string) (*model.User, error) {
	ret := _m.Called(ctx, email, phoneNumber)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmailOrPhoneNumber")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*model.User, error)); ok {
		return rf(ctx, email, phoneNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.User); ok {
		r0 = rf(ctx, email, phoneNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, email, phoneNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserPackage provides a mock function with given fields: ctx, userUID
func (_m *UserUsecase) GetUserPackage(ctx context.Context, userUID string) (*model.UserPackage, error) {
	ret := _m.Called(ctx, userUID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserPackage")
	}

	var r0 *model.UserPackage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.UserPackage, error)); ok {
		return rf(ctx, userUID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.UserPackage); ok {
		r0 = rf(ctx, userUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserPackage)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserUsecase creates a new instance of UserUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserUsecase {
	mock := &UserUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
