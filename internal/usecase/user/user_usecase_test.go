package userusecase_test

import (
	"context"
	"errors"
	"testing"

	"date-apps-be/internal/model"
	"date-apps-be/internal/test"
	userusecase "date-apps-be/internal/usecase/user"
	"date-apps-be/internal/usecase/user/dto"
	"date-apps-be/pkg/datatype"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type params struct {
	Result            *model.User
	UserPackageResult *model.UserPackage
	CreateUser        *dto.CreateUser
	UserUID           string
	Email             string
	PhoneNumber       string
}

func TestCreateUser(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := userusecase.NewUserUsecase(mc.UserRepository, mc.AuthService, mc.UserPremiumRepository)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(token string, err error)
	}{
		{
			caseName: "CreateUser_Success",
			params: params{
				CreateUser: &dto.CreateUser{
					Name:        "John Doe",
					Email:       "john@example.com",
					PhoneNumber: "1234567890",
					Password:    "password123",
				},
			},
			expectations: func(params params) {
				mc.UserRepository.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
					Return(int64(1), nil)
				mc.AuthService.On("GenerateToken", mock.Anything).
					Return("test_token", nil)
			},
			results: func(token string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "test_token", token)
				mc.UserRepository.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything, mock.Anything)
				mc.AuthService.AssertCalled(t, "GenerateToken", mock.Anything)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			token, err := testUsecase.CreateUser(ctx, testCase.params.CreateUser)
			testCase.results(token, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := userusecase.NewUserUsecase(mc.UserRepository, mc.AuthService, mc.UserPremiumRepository)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(user *model.User, err error)
	}{
		{
			caseName: "GetUser_Success",
			params: params{
				UserUID: "test_uid",
				Result: &model.User{
					UID:         "test_uid",
					Name:        "John Doe",
					Email:       ptr("john@example.com"),
					PhoneNumber: ptr("1234567890"),
				},
			},
			expectations: func(params params) {
				mc.UserRepository.On("GetUserByUID", mock.Anything, params.UserUID).Return(params.Result, nil)
			},
			results: func(user *model.User, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, user)
			},
		},
		{
			caseName: "GetUser_NotFound",
			params: params{
				UserUID: "unknown_uid",
			},
			expectations: func(params params) {
				mc.UserRepository.On("GetUserByUID", mock.Anything, "unknown_uid").Return(nil, errors.New("user not found"))
			},
			results: func(user *model.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				mc.UserRepository.AssertCalled(t, "GetUserByUID", mock.Anything, "unknown_uid")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			user, err := testUsecase.GetUser(ctx, testCase.params.UserUID)
			testCase.results(user, err)
		})
	}
}

func TestGetUserByEmailOrPhoneNumber(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := userusecase.NewUserUsecase(mc.UserRepository, mc.AuthService, mc.UserPremiumRepository)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(user *model.User, err error)
	}{
		{
			caseName: "GetUserByEmailOrPhoneNumber_Success",
			params: params{
				Email:       "john@example.com",
				PhoneNumber: "1234567890",
				Result: &model.User{
					UID:         "test_uid",
					Name:        "John Doe",
					Email:       ptr("john@example.com"),
					PhoneNumber: ptr("1234567890"),
				},
			},
			expectations: func(params params) {
				mc.UserRepository.On("GetUserByEmailOrPhoneNumber", mock.Anything, params.Email, params.PhoneNumber).Return(params.Result, nil)
			},
			results: func(user *model.User, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, user)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			user, err := testUsecase.GetUserByEmailOrPhoneNumber(ctx, testCase.params.Email, testCase.params.PhoneNumber)
			testCase.results(user, err)
		})
	}
}

func TestGetUserPackage(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := userusecase.NewUserUsecase(mc.UserRepository, mc.AuthService, mc.UserPremiumRepository)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(userPackage *model.UserPackage, err error)
	}{
		{
			caseName: "GetUserPackage_Success",
			params: params{
				UserPackageResult: &model.UserPackage{
					UserUID:   "2pVJ7Ozfr8obYRW4Dj7bH51BZKU",
					StartedAt: &datatype.Date{},
					EndedAt:   nil,
					Quota:     0,
					PremiumConfig: &model.PremiumConfig{
						UID:         "premium123",
						Name:        "Premium Plan",
						Description: "Premium subscription plan",
						Price:       300,
						Quota:       0,
						ExpiredDay:  0,
						IsActive:    false,
					},
				},
			},
			expectations: func(params params) {
				mc.UserPremiumRepository.On("GetUserPackage", mock.Anything, params.UserUID).Return(params.UserPackageResult, nil)
			},
			results: func(userPackage *model.UserPackage, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, userPackage)
			},
		},
		{
			caseName: "GetUserPackage_NotFound",
			params: params{
				UserUID: "unknown_uid",
			},
			expectations: func(params params) {
				mc.UserPremiumRepository.On("GetUserPackage", mock.Anything, "unknown_uid").Return(nil, errors.New("user package not found"))
			},
			results: func(userPackage *model.UserPackage, err error) {
				assert.Error(t, err)
				assert.Nil(t, userPackage)
				mc.UserPremiumRepository.AssertCalled(t, "GetUserPackage", mock.Anything, "unknown_uid")
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			userPackage, err := testUsecase.GetUserPackage(ctx, testCase.params.UserUID)
			testCase.results(userPackage, err)
		})
	}
}

func ptr(s string) *string {
	return &s
}
