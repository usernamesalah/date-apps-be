package userusecase_test

import (
	"context"
	"errors"
	"testing"

	"date-apps-be/internal/model"
	userrepo "date-apps-be/internal/test/mockrepository"
	authservice "date-apps-be/internal/test/mockservice"
	userusecase "date-apps-be/internal/usecase/user"
	"date-apps-be/internal/usecase/user/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type params struct {
	Result      *model.User
	CreateUser  *dto.CreateUser
	UserUID     string
	Email       string
	PhoneNumber string
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(userrepo.UserRepository)
	mockAuthService := new(authservice.AuthService)
	testUsecase := userusecase.NewUserUsecase(mockUserRepo, mockAuthService)

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
				mockUserRepo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
					Return(int64(1), nil)
				mockAuthService.On("GenerateToken", mock.Anything).
					Return("test_token", nil)
			},
			results: func(token string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "test_token", token)
				mockUserRepo.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything, mock.Anything)
				mockAuthService.AssertCalled(t, "GenerateToken", mock.Anything)
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
	ctx := context.Background()
	mockUserRepo := new(userrepo.UserRepository)
	mockAuthService := new(authservice.AuthService)
	testUsecase := userusecase.NewUserUsecase(mockUserRepo, mockAuthService)

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
					Email:       "john@example.com",
					PhoneNumber: "1234567890",
				},
			},
			expectations: func(params params) {
				mockUserRepo.On("GetUserByUID", mock.Anything, params.UserUID).Return(params.Result, nil)
			},
			results: func(user *model.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, user, user)
				mockUserRepo.AssertCalled(t, "GetUserByUID", mock.Anything, "test_uid")
			},
		},
		{
			caseName: "GetUser_NotFound",
			params: params{
				UserUID: "unknown_uid",
			},
			expectations: func(params params) {
				mockUserRepo.On("GetUserByUID", mock.Anything, "unknown_uid").Return(nil, errors.New("user not found"))
			},
			results: func(user *model.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				mockUserRepo.AssertCalled(t, "GetUserByUID", mock.Anything, "unknown_uid")
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
	ctx := context.Background()
	mockUserRepo := new(userrepo.UserRepository)
	mockAuthService := new(authservice.AuthService)
	testUsecase := userusecase.NewUserUsecase(mockUserRepo, mockAuthService)

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
					Email:       "john@example.com",
					PhoneNumber: "1234567890",
				},
			},
			expectations: func(params params) {
				mockUserRepo.On("GetUserByEmailOrPhoneNumber", mock.Anything, params.Email, params.PhoneNumber).Return(params.Result, nil)
			},
			results: func(user *model.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, user, user)
				mockUserRepo.AssertCalled(t, "GetUserByEmailOrPhoneNumber", mock.Anything, "john@example.com", "1234567890")
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
