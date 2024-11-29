package authservice_test

import (
	"testing"

	"date-apps-be/internal/test"
	"date-apps-be/pkg/derrors"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	mc := test.InitMockComponent(t)
	testAuthService := mc.AuthService

	var testCases = []struct {
		caseName     string
		uid          string
		expectations func()
		results      func(token string, err error)
	}{
		{
			caseName: "GenerateToken_Success",
			uid:      "test_uid",
			expectations: func() {
				mc.AuthService.On("GenerateToken", "test_uid").Return("test_token", nil).Once()
			},
			results: func(token string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "test_token", token)
			},
		},
		{
			caseName: "GenerateToken_Error",
			uid:      "test_uid",
			expectations: func() {
				mc.AuthService.On("GenerateToken", "test_uid").Return("", derrors.New(derrors.Unknown, "token generation error")).Once()
			},
			results: func(token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			// Reset mock before each test case
			mc.AuthService.ExpectedCalls = nil

			testCase.expectations()
			token, err := testAuthService.GenerateToken(testCase.uid)
			testCase.results(token, err)

			mc.AuthService.AssertExpectations(t)
		})
	}
}
