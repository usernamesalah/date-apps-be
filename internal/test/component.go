package test

import (
	"crypto/rand"
	"crypto/rsa"
	"date-apps-be/infrastructure/config"
	"date-apps-be/internal/test/mockrepository"
	"date-apps-be/internal/test/mockservice"
	"date-apps-be/internal/test/mockusecase"
	"testing"
)

type MockComponent struct {
	Config                  *config.Config
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
		Config: &config.Config{
			JWTRS256PrivateKey: generatePrivateKey(),
			JWTExpiration:      10,
		},
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

func generatePrivateKey() *rsa.PrivateKey {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privKey
}
