package handler

import (
	"date-apps-be/internal/api/http/handler/request"
	"date-apps-be/internal/container"
	"date-apps-be/internal/model"
	authservice "date-apps-be/internal/service/auth"
	userusecase "date-apps-be/internal/usecase/user"
	"date-apps-be/internal/usecase/user/dto"
	"date-apps-be/pkg/api"
	"date-apps-be/pkg/derrors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	userHandler struct {
		userUsecase userusecase.UserUsecase
		authservice authservice.AuthService
	}

	UserHandler interface {
		GetUserProfile(c echo.Context) error
		Login(c echo.Context) error
		Register(c echo.Context) error
	}
)

func NewUserHandler(hc *container.HandlerComponent) UserHandler {
	return &userHandler{
		userUsecase: hc.UserUsecase,
		authservice: hc.AuthService,
	}
}

// GetUserProfile retrieves the user's profile information.
// It first retrieves the user's UID from the context.
func (u *userHandler) GetUserProfile(c echo.Context) error {
	userInfo := c.Get("userInfo").(*model.JWTClaims)

	user, err := u.userUsecase.GetUser(c.Request().Context(), userInfo.UserUID)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, user, http.StatusOK)
}

// RegisterUser creates a new user account.
// It reads the request body and decodes it into a User model.
func (u *userHandler) Register(c echo.Context) error {
	req := new(request.UserRegister)
	if err := c.Bind(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if err := c.Validate(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	userExist, err := u.userUsecase.GetUserByEmailOrPhoneNumber(c.Request().Context(), req.Email, req.PhoneNumber)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if userExist != nil {
		return api.RenderErrorResponse(c, c.Request(), derrors.New(derrors.InvalidArgument, "Email or Phone Number already registered, please login"))
	}

	t, err := u.userUsecase.CreateUser(c.Request().Context(), &dto.CreateUser{
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
	})
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, map[string]string{
		"token": t,
	}, http.StatusOK)
}

// Login
func (u *userHandler) Login(c echo.Context) error {
	req := new(request.UserLogin)
	if err := c.Bind(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if err := c.Validate(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	user, err := u.userUsecase.GetUserByEmailOrPhoneNumber(c.Request().Context(), req.Email, req.PhoneNumber)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if user == nil {
		return api.RenderErrorResponse(c, c.Request(), derrors.New(derrors.NotFound, "Email or Phone Number not found, please register"))
	}

	token, err := u.authservice.GenerateToken(user.UID)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, map[string]string{
		"token": token,
	}, http.StatusOK)
}
