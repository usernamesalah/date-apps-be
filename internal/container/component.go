package container

import (
	"date-apps-be/infrastructure/config"
	repository "date-apps-be/internal/repository/common"
	userrepository "date-apps-be/internal/repository/user"
	usermatchrepository "date-apps-be/internal/repository/user_match"
	authservice "date-apps-be/internal/service/auth"
	userusecase "date-apps-be/internal/usecase/user"
	usermatchusecase "date-apps-be/internal/usecase/user_match"
)

type HandlerComponent struct {
	Config *config.Config

	// Service
	AuthService authservice.AuthService

	// Usecase
	UserUsecase      userusecase.UserUsecase
	UserMatchUsecase usermatchusecase.UserMatchUsecase
}

func NewHandlerComponent(sc *SharedComponent) *HandlerComponent {

	baseStore := repository.NewRepository(sc.DB)

	authservice := authservice.NewAuthService(sc.Conf)

	userRepo := userrepository.NewUserRepository(baseStore)
	userUsecase := userusecase.NewUserUsecase(userRepo, authservice)

	userMatchRepo := usermatchrepository.NewUserMatchRepository(baseStore)
	userMatchUsecase := usermatchusecase.NewUserMatchUsecase(userMatchRepo)

	return &HandlerComponent{
		Config: sc.Conf,
		// Service
		AuthService: authservice,

		// Usecase
		UserUsecase:      userUsecase,
		UserMatchUsecase: userMatchUsecase,
	}
}
