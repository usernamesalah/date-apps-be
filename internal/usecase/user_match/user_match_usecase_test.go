package usermatchusecase_test

import (
	"context"
	"errors"
	"testing"

	"date-apps-be/internal/constant"
	"date-apps-be/internal/model"
	"date-apps-be/internal/test"
	usermatchusecase "date-apps-be/internal/usecase/user_match"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type params struct {
	UserMatch   *model.UserMatch
	UserPackage *model.UserPackage
	UserUID     string
	MatchUID    string
}

func TestCreateUserMatch(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := usermatchusecase.NewUserMatchUsecase(mc.UserMatchRepository, mc.UserUsecase)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(err error)
	}{
		{
			caseName: "CreateUserMatch_Success",
			params: params{
				UserMatch: &model.UserMatch{
					UserUID:  "user123",
					MatchUID: "match123",
				},
				UserPackage: &model.UserPackage{
					UserUID: "user123",
					Quota:   5,
				},
			},
			expectations: func(params params) {
				mc.UserUsecase.On("GetUserPackage", mock.Anything, params.UserMatch.UserUID).Return(params.UserPackage, nil)
				mc.UserMatchRepository.On("GetTotalUserMatchToday", mock.Anything, params.UserMatch.UserUID).Return(0, nil)
				mc.UserMatchRepository.On("CreateUserMatch", mock.Anything, mock.Anything).Return(nil)
			},
			results: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			caseName: "CreateUserMatch_ExceededQuota",
			params: params{
				UserMatch: &model.UserMatch{
					UserUID:  "user123",
					MatchUID: "match123",
				},
				UserPackage: &model.UserPackage{
					UserUID: "user123",
					Quota:   0,
				},
			},
			expectations: func(params params) {
				mc.UserUsecase.On("GetUserPackage", mock.Anything, params.UserMatch.UserUID).Return(params.UserPackage, nil)
				mc.UserMatchRepository.On("GetTotalUserMatchToday", mock.Anything, params.UserMatch.UserUID).Return(constant.MaxMatchPerDay, nil)
			},
			results: func(err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			err := testUsecase.CreateUserMatch(ctx, testCase.params.UserMatch)
			testCase.results(err)
		})
	}
}

func TestGetUserMatchTodayByUserUIDAndMatchUID(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := usermatchusecase.NewUserMatchUsecase(mc.UserMatchRepository, mc.UserUsecase)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(userMatch *model.UserMatch, err error)
	}{
		{
			caseName: "GetUserMatchTodayByUserUIDAndMatchUID_Success",
			params: params{
				UserUID:  "user123",
				MatchUID: "match123",
				UserMatch: &model.UserMatch{
					UserUID:  "user123",
					MatchUID: "match123",
				},
			},
			expectations: func(params params) {
				mc.UserMatchRepository.On("GetUserMatchTodayByUserUIDAndMatchUID", mock.Anything, params.UserUID, params.MatchUID).Return(params.UserMatch, nil)
			},
			results: func(userMatch *model.UserMatch, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, userMatch)
				mc.UserMatchRepository.AssertCalled(t, "GetUserMatchTodayByUserUIDAndMatchUID", mock.Anything, "user123", "match123")
			},
		},
		{
			caseName: "GetUserMatchTodayByUserUIDAndMatchUID_NotFound",
			params: params{
				UserUID:  "user123",
				MatchUID: "unknown_match",
			},
			expectations: func(params params) {
				mc.UserMatchRepository.On("GetUserMatchTodayByUserUIDAndMatchUID", mock.Anything, params.UserUID, params.MatchUID).Return(nil, errors.New("user match not found"))
			},
			results: func(userMatch *model.UserMatch, err error) {
				assert.Error(t, err)
				assert.Nil(t, userMatch)
				mc.UserMatchRepository.AssertCalled(t, "GetUserMatchTodayByUserUIDAndMatchUID", mock.Anything, "user123", "unknown_match")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			userMatch, err := testUsecase.GetUserMatchTodayByUserUIDAndMatchUID(ctx, testCase.params.UserUID, testCase.params.MatchUID)
			testCase.results(userMatch, err)
		})
	}
}
