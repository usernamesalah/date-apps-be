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
		GetMyPackage(c echo.Context) error
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
// Get user profile
// @Summary Get user profile
// @Description Get user profile
// @Tags users
// @ID get-user-uid
// @Produce json
// @Param authorization header string true "bearer token"
// @Success 200 {object} map[string]string
// @Router /users/profile [get]
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
// RegisterUser creates a new user account.
// It reads the request body and decodes it into a User model.
// Register user
// @Summary Register user
// @Description Register a new user
// @Tags auth
// @ID register-user
// @Accept json
// @Produce json
// @Param user body request.UserRegister true "User registration details"
// @Success 200 {object} map[string]string
// @Router /register [post]
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
// Login user
// @Summary Login user
// @Description Authenticate user and return a JWT token
// @Tags auth
// @ID login-user
// @Accept json
// @Produce json
// @Param user body request.UserLogin true "User login details"
// @Success 200 {object} map[string]string
// @Router /login [post]
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

// GetMyPackage retrieves the user's package information.
// It first retrieves the user's UID from the context.
// GetMyPackage retrieves the user's package information.
// @Summary Get user package
// @Description Get user package information
// @Tags users
// @ID get-user-package
// @Produce json
// @Param authorization header string true "bearer token"
// @Success 200 {object} map[string]string
// @Router /users/package [get]
func (u *userHandler) GetMyPackage(c echo.Context) error {
	userInfo := c.Get("userInfo").(*model.JWTClaims)

	userPackage, err := u.userUsecase.GetUserPackage(c.Request().Context(), userInfo.UserUID)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, userPackage, http.StatusOK)
}
