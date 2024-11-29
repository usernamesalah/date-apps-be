package premiumconfigusecase_test

import (
	"context"
	"errors"
	"testing"

	"date-apps-be/internal/model"
	"date-apps-be/internal/test"
	premiumconfigusecase "date-apps-be/internal/usecase/premium_config"
	"date-apps-be/internal/usecase/premium_config/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type params struct {
	PremiumConfig    *model.PremiumConfig
	PremiumConfigs   []*model.PremiumConfig
	UserPackage      *model.UserPackage
	UserPurchase     dto.UserPurchase
	PremiumConfigUID string
	UserUID          string
	Page             uint64
	Limit            uint64
}

func TestGetPremiumConfigs(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := premiumconfigusecase.NewPremiumConfigUsecase(mc.PremiumConfigRepository, mc.UserPremiumRepository)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(configs []*model.PremiumConfig, err error)
	}{
		{
			caseName: "GetPremiumConfigs_Success",
			params: params{
				Page:  1,
				Limit: 10,
				PremiumConfigs: []*model.PremiumConfig{
					{
						UID:         "premium123",
						Name:        "Premium Plan",
						Description: "Premium subscription plan",
						Price:       300,
						Quota:       10,
						ExpiredDay:  30,
						IsActive:    true,
					},
				},
			},
			expectations: func(params params) {
				mc.PremiumConfigRepository.On("GetPremiumConfigs", mock.Anything, params.Page, params.Limit).Return(params.PremiumConfigs, nil)
			},
			results: func(configs []*model.PremiumConfig, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, configs)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			configs, err := testUsecase.GetPremiumConfigs(ctx, testCase.params.Page, testCase.params.Limit)
			testCase.results(configs, err)
		})
	}
}

func TestGetPremiumConfigByUID(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := premiumconfigusecase.NewPremiumConfigUsecase(mc.PremiumConfigRepository, mc.UserPremiumRepository)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(config *model.PremiumConfig, err error)
	}{
		{
			caseName: "GetPremiumConfigByUID_Success",
			params: params{
				PremiumConfigUID: "premium123",
				PremiumConfig: &model.PremiumConfig{
					UID:         "premium123",
					Name:        "Premium Plan",
					Description: "Premium subscription plan",
					Price:       300,
					Quota:       10,
					ExpiredDay:  30,
					IsActive:    true,
				},
			},
			expectations: func(params params) {
				mc.PremiumConfigRepository.On("GetPremiumConfigByUID", mock.Anything, params.PremiumConfigUID).Return(params.PremiumConfig, nil)
			},
			results: func(config *model.PremiumConfig, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, config)
			},
		},
		{
			caseName: "GetPremiumConfigByUID_NotFound",
			params: params{
				PremiumConfigUID: "unknown_uid",
			},
			expectations: func(params params) {
				mc.PremiumConfigRepository.On("GetPremiumConfigByUID", mock.Anything, params.PremiumConfigUID).Return(nil, errors.New("premium config not found"))
			},
			results: func(config *model.PremiumConfig, err error) {
				assert.Error(t, err)
				assert.Nil(t, config)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			config, err := testUsecase.GetPremiumConfigByUID(ctx, testCase.params.PremiumConfigUID)
			testCase.results(config, err)
		})
	}
}

func TestPurchasePackage(t *testing.T) {
	mc := test.InitMockComponent(t)
	ctx := context.Background()
	testUsecase := premiumconfigusecase.NewPremiumConfigUsecase(mc.PremiumConfigRepository, mc.UserPremiumRepository)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(err error)
	}{
		{
			caseName: "PurchasePackage_Success",
			params: params{
				UserPurchase: dto.UserPurchase{
					UserUID:          "user123",
					PremiumConfigUID: "premium123",
				},
				PremiumConfig: &model.PremiumConfig{
					UID:         "premium123",
					Name:        "Premium Plan",
					Description: "Premium subscription plan",
					Price:       300,
					Quota:       10,
					ExpiredDay:  30,
					IsActive:    true,
				},
			},
			expectations: func(params params) {
				mc.UserPremiumRepository.On("GetUserPackage", mock.Anything, params.UserPurchase.UserUID).Return(nil, nil)
				mc.PremiumConfigRepository.On("GetPremiumConfigByUID", mock.Anything, params.UserPurchase.PremiumConfigUID).Return(params.PremiumConfig, nil)
				// Updated context matching here
				mc.UserPremiumRepository.On("CreateUserPackage", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), mock.Anything, mock.AnythingOfType("*model.UserPackage")).Return(nil)
			},
			results: func(err error) {
				assert.NoError(t, err)
				mc.UserPremiumRepository.AssertCalled(t, "GetUserPackage", mock.Anything, "user123")
				mc.PremiumConfigRepository.AssertCalled(t, "GetPremiumConfigByUID", mock.Anything, "premium123")
				mc.UserPremiumRepository.AssertCalled(t, "CreateUserPackage", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), mock.Anything, mock.AnythingOfType("*model.UserPackage"))
			},
		},
		{
			caseName: "PurchasePackage_UserAlreadyHasPackage",
			params: params{
				UserPurchase: dto.UserPurchase{
					UserUID:          "user123",
					PremiumConfigUID: "premium123",
				},
				UserPackage: &model.UserPackage{
					UserUID: "user123",
				},
			},
			expectations: func(params params) {
				mc.UserPremiumRepository.On("GetUserPackage", mock.Anything, params.UserPurchase.UserUID).Return(params.UserPackage, nil)
			},
			results: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			err := testUsecase.PurchasePackage(ctx, testCase.params.UserPurchase)
			testCase.results(err)
		})
	}
}
