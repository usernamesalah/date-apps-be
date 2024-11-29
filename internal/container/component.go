package container

import (
	"date-apps-be/infrastructure/config"
	repository "date-apps-be/internal/repository/common"
	premiumconfigrepository "date-apps-be/internal/repository/premium_config"
	userrepository "date-apps-be/internal/repository/user"
	usermatchrepository "date-apps-be/internal/repository/user_match"
	userpackagerepository "date-apps-be/internal/repository/user_premium"
	authservice "date-apps-be/internal/service/auth"
	premiumconfigusecase "date-apps-be/internal/usecase/premium_config"
	userusecase "date-apps-be/internal/usecase/user"
	usermatchusecase "date-apps-be/internal/usecase/user_match"
)

type HandlerComponent struct {
	Config *config.Config

	// Service
	AuthService authservice.AuthService

	// Usecase
	UserUsecase          userusecase.UserUsecase
	UserMatchUsecase     usermatchusecase.UserMatchUsecase
	PremiumConfigUsecase premiumconfigusecase.PremiumConfigUsecase
}

func NewHandlerComponent(sc *SharedComponent) *HandlerComponent {

	baseStore := repository.NewRepository(sc.DB)

	authservice := authservice.NewAuthService(sc.Conf)

	userPackageRepo := userpackagerepository.NewUserPremiumRepository(baseStore)
	userRepo := userrepository.NewUserRepository(baseStore)
	userUsecase := userusecase.NewUserUsecase(userRepo, authservice, userPackageRepo)

	userMatchRepo := usermatchrepository.NewUserMatchRepository(baseStore)
	userMatchUsecase := usermatchusecase.NewUserMatchUsecase(userMatchRepo, userUsecase)

	premiumConfigRepo := premiumconfigrepository.NewPremiumConfigRepository(baseStore)
	premiumConfigUsecase := premiumconfigusecase.NewPremiumConfigUsecase(premiumConfigRepo, userPackageRepo)

	return &HandlerComponent{
		Config: sc.Conf,

		// Service
		AuthService: authservice,

		// Usecase
		UserUsecase:          userUsecase,
		UserMatchUsecase:     userMatchUsecase,
		PremiumConfigUsecase: premiumConfigUsecase,
	}
}
