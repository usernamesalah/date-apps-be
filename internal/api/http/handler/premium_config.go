package handler

import (
	"date-apps-be/internal/api/http/handler/request"
	"date-apps-be/internal/container"
	"date-apps-be/internal/model"
	premiumconfigusecase "date-apps-be/internal/usecase/premium_config"
	"date-apps-be/internal/usecase/premium_config/dto"
	"date-apps-be/pkg/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	premiumConfigHandler struct {
		premiumConfigUsecase premiumconfigusecase.PremiumConfigUsecase
	}

	PremiumConfigHandler interface {
		GetPackages(c echo.Context) error
		GetPackageByUID(c echo.Context) error
		PurchasePackage(c echo.Context) error
	}
)

func NewPremiumConfigHandler(hc *container.HandlerComponent) PremiumConfigHandler {
	return &premiumConfigHandler{
		premiumConfigUsecase: hc.PremiumConfigUsecase,
	}
}

// GetPackages retrieves a list of available premium packages.
// It retrieves the pagination parameters and retrieves the list of available packages.
// @Summary Get available premium packages
// @Description Retrieves a list of available premium packages with pagination
// @Tags premium
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} []model.PremiumConfig "List of premium packages"
// @Router /packages [get]
func (p *premiumConfigHandler) GetPackages(c echo.Context) error {
	page, limit, err := api.ParsePagination(c.Request())
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	configs, err := p.premiumConfigUsecase.GetPremiumConfigs(c.Request().Context(), page, limit)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, configs, http.StatusOK)
}

// GetPackageByUID retrieves a premium package by its UID.
// It retrieves the package UID from the request parameters and retrieves the package details.
// @Summary Get premium package by UID
// @Description Retrieves a premium package by its UID
// @Tags premium
// @Accept json
// @Produce json
// @Param uid path string true "Package UID"
// @Success 200 {object} model.PremiumConfig "Premium package details"
// @Router /packages/{uid} [get]
func (p *premiumConfigHandler) GetPackageByUID(c echo.Context) error {
	uid := c.Param("uid")

	config, err := p.premiumConfigUsecase.GetPremiumConfigByUID(c.Request().Context(), uid)
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseOK(c, config, http.StatusOK)
}

// PurchasePackage purchases a premium package for the user.
// It retrieves the user's UID from the context and the package UID from the request body.
// @Summary Purchase a premium package
// @Description Purchases a premium package for the user
// @Tags premium
// @Accept json
// @Produce json
// @Param authorization header string true "bearer token"
// @Param userPurchase body request.UserPurchase true "User purchase request"
// @Success 200 {object} map[string]string "Successfully purchased the package"
// @Router /packages/purchase [post]
func (p *premiumConfigHandler) PurchasePackage(c echo.Context) error {
	userInfo := c.Get("userInfo").(*model.JWTClaims)

	req := new(request.UserPurchase)
	if err := c.Bind(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	if err := c.Validate(req); err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	err := p.premiumConfigUsecase.PurchasePackage(c.Request().Context(), dto.UserPurchase{
		UserUID:          userInfo.UserUID,
		PremiumConfigUID: req.PremiumConfigUID,
	})
	if err != nil {
		return api.RenderErrorResponse(c, c.Request(), err)
	}

	return api.ResponseSuccess(c, nil, "success Purchase the package", http.StatusOK)
}
