package handler_test

import (
	"date-apps-be/internal/api/http/handler"
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

func ptr(s string) *string {
	return &s
}
func TestUserHandler_GetUserProfile(t *testing.T) {
	e := echo.New()
	mockComponent := test.InitMockComponent(t)

	hc := &container.HandlerComponent{
		UserUsecase: mockComponent.UserUsecase,
		AuthService: mockComponent.AuthService,
	}

	h := handler.NewUserHandler(hc)

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedUser   *model.User
	}{
		{
			name: "success get user profile",
			setupMock: func() {
				mockComponent.UserUsecase.On("GetUser",
					mock.Anything,
					"test-uid",
				).Return(&model.User{
					UID:         "test-uid",
					Name:        "Test User",
					Email:       ptr("test@example.com"),
					PhoneNumber: ptr("1234567890"),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedUser: &model.User{
				UID:         "test-uid",
				Name:        "Test User",
				Email:       ptr("test@example.com"),
				PhoneNumber: ptr("1234567890"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set user info in context
			c.Set("userInfo", &model.JWTClaims{UserUID: "test-uid"})

			err := h.GetUserProfile(c)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var response struct {
				Data *model.User `json:"data"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUser.UID, response.Data.UID)
			assert.Equal(t, tc.expectedUser.Name, response.Data.Name)
			assert.Equal(t, tc.expectedUser.Email, response.Data.Email)
			assert.Equal(t, tc.expectedUser.PhoneNumber, response.Data.PhoneNumber)
		})
	}
}

func TestUserHandler_Register(t *testing.T) {
	e := echo.New()
	e.Validator = &requestValidator{}
	mockComponent := test.InitMockComponent(t)

	hc := &container.HandlerComponent{
		UserUsecase: mockComponent.UserUsecase,
		AuthService: mockComponent.AuthService,
	}

	h := handler.NewUserHandler(hc)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedError  error
	}{
		{
			name: "success register",
			requestBody: `{
				"email": "test@example.com",
				"password": "password123",
				"phone_number": "1234567890",
				"name": "Test User"
			}`,
			setupMock: func() {
				mockComponent.UserUsecase.On("GetUserByEmailOrPhoneNumber",
					mock.Anything,
					"test@example.com",
					"1234567890",
				).Return(nil, nil)

				mockComponent.UserUsecase.On("CreateUser",
					mock.Anything,
					mock.AnythingOfType("*dto.CreateUser"),
				).Return("jwt-token", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "user already exists",
			requestBody: `{
				"email": "existing@example.com",
				"password": "password123",
				"phone_number": "1234567890",
				"name": "Test User"
			}`,
			setupMock: func() {
				mockComponent.UserUsecase.On("GetUserByEmailOrPhoneNumber",
					mock.Anything,
					"existing@example.com",
					"1234567890",
				).Return(&model.User{
					UID:         "existing-uid",
					Name:        "Existing User",
					Email:       ptr("existing@example.com"),
					PhoneNumber: ptr("1234567890"),
				}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.Register(c)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
