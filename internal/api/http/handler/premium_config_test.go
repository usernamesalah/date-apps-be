package handler_test

import (
	"date-apps-be/internal/api/http/handler"
	"date-apps-be/internal/container"
	"date-apps-be/internal/model"
	"date-apps-be/internal/test"
	"date-apps-be/pkg/derrors"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPremiumConfigHandler_GetPackages(t *testing.T) {
	e := echo.New()
	e.Validator = NewValidator()
	mockComponent := test.InitMockComponent(t)

	hc := &container.HandlerComponent{
		PremiumConfigUsecase: mockComponent.PremiumConfigUsecase,
	}

	h := handler.NewPremiumConfigHandler(hc)

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedData   []*model.PremiumConfig
		expectedError  *derrors.Error
	}{
		{
			name: "success get packages",
			setupMock: func() {
				mockComponent.PremiumConfigUsecase.On("GetPremiumConfigs",
					mock.Anything,
					uint64(1),
					uint64(10),
				).Return([]*model.PremiumConfig{
					{
						UID:        "premium-1",
						Name:       "Premium Package",
						Price:      100,
						ExpiredDay: 30,
						IsActive:   true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedData: []*model.PremiumConfig{
				{
					UID:        "premium-1",
					Name:       "Premium Package",
					Price:      100,
					ExpiredDay: 30,
					IsActive:   true,
				},
			},
		},
		{
			name: "internal server error",
			setupMock: func() {
				mockComponent.PremiumConfigUsecase.On("GetPremiumConfigs",
					mock.Anything,
					uint64(1),
					uint64(10),
				).Return(nil, derrors.New(derrors.Unknown, "OK"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  derrors.New(derrors.Unknown, "OK").(*derrors.Error),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/packages?page=1&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.GetPackages(c)
			assert.NoError(t, err)

			if tc.expectedError != nil {
				var errorResponse struct {
					Message string `json:"message"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Contains(t, errorResponse.Message, tc.expectedError.Error())
			} else {
				var response struct {
					Data []*model.PremiumConfig `json:"data"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, response.Data)
			}
		})
	}
}

func NewValidator() *requestValidator {
	return &requestValidator{}
}

type requestValidator struct{}

func (rv *requestValidator) Validate(i interface{}) (err error) {
	_, err = govalidator.ValidateStruct(i)
	return
}
