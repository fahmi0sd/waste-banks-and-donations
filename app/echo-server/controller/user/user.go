package user

import (
	"log/slog"
	"net/http"

	"github.com/fahmi0sd/waste-banks-and-donations/service/user"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	service  user.Service
	validate *validator.Validate
}

func NewController(logger *slog.Logger, service user.Service) *Controller {
	return &Controller{logger: logger, service: service, validate: validator.New()}
}

func (ctrl *Controller) Me(c echo.Context) error {
	userID := auth.UserID(c)

	profile, err := ctrl.service.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("profil ditemukan", profile))
}

func (ctrl *Controller) UpdateMe(c echo.Context) error {
	userID := auth.UserID(c)

	var req user.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("data yang dikirim tidak valid"))
	}

	profile, err := ctrl.service.UpdateProfile(userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("profil berhasil diperbarui", profile))
}
