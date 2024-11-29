package userusecase

import (
	"context"
	"date-apps-be/internal/model"
	userrepo "date-apps-be/internal/repository/user"
	userpackagerepo "date-apps-be/internal/repository/user_premium"
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
		GetUserPackage(ctx context.Context, userUID string) (userPackage *model.UserPackage, err error)
	}

	userUsecase struct {
		authService authservice.AuthService
		userRepo    userrepo.UserRepository
		userPackage userpackagerepo.UserPremiumRepository
	}
)

func NewUserUsecase(userRepo userrepo.UserRepository, authService authservice.AuthService, userPackage userpackagerepo.UserPremiumRepository) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		authService: authService,
		userPackage: userPackage,
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
		Email:       &user.Email,
		PhoneNumber: &user.PhoneNumber,
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
	if err != nil {
		return
	}

	userPackage, err := u.userPackage.GetUserPackage(ctx, userUID)
	if err != nil {
		return
	}

	if userPackage != nil {
		user.IsPremium = true
	}
	return
}

func (u *userUsecase) GetUserByEmailOrPhoneNumber(ctx context.Context, email, phoneNumber string) (user *model.User, err error) {
	defer derrors.Wrap(&err, "GetUserByEmailOrPhoneNumber(%q , %q)", email, phoneNumber)

	user, err = u.userRepo.GetUserByEmailOrPhoneNumber(ctx, email, phoneNumber)
	return
}

func (u *userUsecase) GetUserPackage(ctx context.Context, userUID string) (userPackage *model.UserPackage, err error) {
	defer derrors.Wrap(&err, "GetUserPackage(%q)", userUID)
	userPackage, err = u.userPackage.GetUserPackage(ctx, userUID)
	return
}
