package test

import (
	"date-apps-be/internal/test/mockrepository"
	"date-apps-be/internal/test/mockservice"
	"date-apps-be/internal/test/mockusecase"
	"testing"
)

type MockComponent struct {
	UserRepository          *mockrepository.UserRepository
	UserMatchRepository     *mockrepository.UserMatchRepository
	UserPremiumRepository   *mockrepository.UserPremiumRepository
	PremiumConfigRepository *mockrepository.PremiumConfigRepository
	UserUsecase             *mockusecase.UserUsecase
	UserMatchUsecase        *mockusecase.UserMatchUsecase
	PremiumConfigUsecase    *mockusecase.PremiumConfigUsecase
	AuthService             *mockservice.AuthService
}

func InitMockComponent(t *testing.T) *MockComponent {
	return &MockComponent{
		UserRepository:          mockrepository.NewUserRepository(t),
		UserMatchRepository:     mockrepository.NewUserMatchRepository(t),
		UserPremiumRepository:   mockrepository.NewUserPremiumRepository(t),
		PremiumConfigRepository: mockrepository.NewPremiumConfigRepository(t),
		UserUsecase:             mockusecase.NewUserUsecase(t),
		UserMatchUsecase:        mockusecase.NewUserMatchUsecase(t),
		PremiumConfigUsecase:    mockusecase.NewPremiumConfigUsecase(t),
		AuthService:             mockservice.NewAuthService(t),
	}
}
