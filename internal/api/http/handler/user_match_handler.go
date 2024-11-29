package handler

import (
	"date-apps-be/internal/api/http/handler/request"
	"date-apps-be/internal/api/http/handler/response"
	"date-apps-be/internal/constant"
	"date-apps-be/internal/container"
	"date-apps-be/internal/model"
	userMatchUsecase "date-apps-be/internal/usecase/user_match"
	"date-apps-be/pkg/api"
	"date-apps-be/pkg/derrors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserMatchHandler defines the interface for handling user match-related HTTP requests.
// It includes methods for creating a match and retrieving user matches.
type (
	UserMatchHandler interface {
		CreateMatch(c echo.Context) error
		GetUserMatches(c echo.Context) error
	}

	userMatchHandler struct {
		userMatchUsecase userMatchUsecase.UserMatchUsecase
	}
)

func NewUserMatchHandler(hc *container.HandlerComponent) UserMatchHandler {
	return &userMatchHandler{
		userMatchUsecase: hc.UserMatchUsecase,
	}
}

// GetUserMatches retrieves a list of available users for the current user.
// @Summary Get user matches
// @Accept json
// @Tags UserMatch
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} response.UserMatchResponse "List of available users and remaining quota"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Param authorization header string true "bearer token"
// @Router /matches [get]
func (u *userMatchHandler) GetUserMatches(c echo.Context) error {
	userInfo := c.Get("userInfo").(*model.JWTClaims)

	page, limit, err := api.ParsePagination(c.Request())
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	users, quotaLeft, err := u.userMatchUsecase.GetAvailableUsers(c.Request().Context(), userInfo.UserUID, page, limit)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, response.NewUserMatchResponse(users, quotaLeft), http.StatusOK)
}

// CreateMatch handles the creation of a user match.
// It retrieves user information from the context, binds and validates the request,
// parses the match type, checks if a match already exists for the user today,
// and if not, creates a new user match.
// It returns an appropriate response based on the success or failure of these operations.
// @Summary Create a user match
// @Tags UserMatch
// @Accept json
// @Produce json
// @Param authorization header string true "bearer token"
// @Param req body request.CreateMatch true "Create Match Request"
// @Success 201 {object} map[string]string "Success Match with that Person"
// @Router /matches [post]
func (u *userMatchHandler) CreateMatch(c echo.Context) error {
	userInfo := c.Get("userInfo").(*model.JWTClaims)

	req := new(request.CreateMatch)
	if err := c.Bind(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if err := c.Validate(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	matchType, err := constant.ParseUserMatchType(req.MatchType)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	userMatchExist, err := u.userMatchUsecase.GetUserMatchTodayByUserUIDAndMatchUID(c.Request().Context(), userInfo.UserUID, req.MatchUID)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if userMatchExist != nil {
		return api.RenderErrorResponse(c, c.Request(), derrors.New(derrors.Forbidden, "you already matched with this user today"))
	}

	userMatch := model.UserMatch{
		UserUID:   userInfo.UserUID,
		MatchUID:  req.MatchUID,
		MatchType: matchType,
	}

	err = u.userMatchUsecase.CreateUserMatch(c.Request().Context(), &userMatch)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseSuccess(c, nil, "Success Match with that Person", http.StatusCreated)
}
