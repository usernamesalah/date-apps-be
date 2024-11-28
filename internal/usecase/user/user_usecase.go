package userusecase

import (
	"context"
	"date-apps-be/internal/model"
	userrepo "date-apps-be/internal/repository/user"
	authservice "date-apps-be/internal/service/auth"
	"date-apps-be/internal/usecase/user/dto"
	"date-apps-be/pkg/derrors"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserUsecase interface {
		CreateUser(ctx context.Context, user *dto.CreateUser) (token string, err error)
		GetUser(ctx context.Context, userUID string) (user *model.User, err error)
		GetUserByEmailOrPhoneNumber(ctx context.Context, email, phoneNumber string) (user *model.User, err error)
	}

	userUsecase struct {
		userRepo    userrepo.UserRepository
		authService authservice.AuthService
	}
)

func NewUserUsecase(userRepo userrepo.UserRepository, authService authservice.AuthService) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *dto.CreateUser) (token string, err error) {
	defer derrors.Wrap(&err, "CreateUser(%q)", user.Name)

	uid := ksuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	userData := &model.User{
		UID:         uid,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Password:    string(hashedPassword),
	}

	_, err = u.userRepo.CreateUser(ctx, nil, userData)
	if err != nil {
		return
	}

	token, err = u.authService.GenerateToken(uid)
	if err != nil {
		return
	}

	return
}

func (u *userUsecase) GetUser(ctx context.Context, userUID string) (user *model.User, err error) {
	defer derrors.Wrap(&err, "GetUser(%q)", userUID)

	user, err = u.userRepo.GetUserByUID(ctx, userUID)
	return
}

func (u *userUsecase) GetUserByEmailOrPhoneNumber(ctx context.Context, email, phoneNumber string) (user *model.User, err error) {
	defer derrors.Wrap(&err, "GetUserByEmailOrPhoneNumber(%q , %q)", email, phoneNumber)

	user, err = u.userRepo.GetUserByEmailOrPhoneNumber(ctx, email, phoneNumber)
	return
}
