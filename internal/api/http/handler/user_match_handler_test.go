package handler_test

import (
	"date-apps-be/internal/api/http/handler"
	"date-apps-be/internal/api/http/handler/response"
	"date-apps-be/internal/constant"
	"date-apps-be/internal/container"
	"date-apps-be/internal/model"
	"date-apps-be/internal/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserMatchHandler_GetUserMatches(t *testing.T) {
	// Setup
	e := echo.New()
	mockComponent := test.InitMockComponent(t)

	hc := &container.HandlerComponent{
		UserMatchUsecase: mockComponent.UserMatchUsecase,
	}

	h := handler.NewUserMatchHandler(hc)

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedUsers  []*model.User
		expectedQuota  int
	}{
		{
			name: "success get user matches",
			setupMock: func() {
				mockComponent.UserMatchUsecase.On("GetAvailableUsers",
					mock.Anything,
					"test-uid",
					uint64(1),
					uint64(10),
				).Return([]*model.User{
					{UID: "user-1", Name: "Test User 1"},
					{UID: "user-2", Name: "Test User 2"},
				}, 5, nil)
			},
			expectedStatus: http.StatusOK,
			expectedUsers: []*model.User{
				{UID: "user-1", Name: "Test User 1"},
				{UID: "user-2", Name: "Test User 2"},
			},
			expectedQuota: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			tc.setupMock()

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/matches?page=1&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set user info in context
			c.Set("userInfo", &model.JWTClaims{UserUID: "test-uid"})

			// Execute request
			err := h.GetUserMatches(c)
			assert.NoError(t, err)

			// Assert response
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var response struct {
				Data response.UserMatchResponse `json:"data"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, len(tc.expectedUsers), len(response.Data.Users))
			assert.Equal(t, tc.expectedQuota, response.Data.QuotaLeft)
		})
	}
}

func TestUserMatchHandler_CreateMatch(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &requestValidator{}
	mockComponent := test.InitMockComponent(t)

	hc := &container.HandlerComponent{
		UserMatchUsecase: mockComponent.UserMatchUsecase,
	}

	h := handler.NewUserMatchHandler(hc)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:        "success create match",
			requestBody: `{"match_uid":"match-uid","match_type":"like"}`,
			setupMock: func() {
				mockComponent.UserMatchUsecase.On("GetUserMatchTodayByUserUIDAndMatchUID",
					mock.Anything,
					"test-uid",
					"match-uid",
				).Return(nil, nil)

				mockComponent.UserMatchUsecase.On("CreateUserMatch",
					mock.Anything,
					mock.MatchedBy(func(match *model.UserMatch) bool {
						return match.UserUID == "test-uid" &&
							match.MatchUID == "match-uid" &&
							match.MatchType == constant.UserMatchTypeLike
					}),
				).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			tc.setupMock()

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/matches",
				strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set user info in context
			c.Set("userInfo", &model.JWTClaims{UserUID: "test-uid"})

			// Execute request
			err := h.CreateMatch(c)
			assert.NoError(t, err)

			// Assert response
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
