package auth

import (
	"log/slog"
	"net/http"

	"github.com/fahmi0sd/waste-banks-and-donations/service/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	service  auth.Service
	validate *validator.Validate
}

func NewController(logger *slog.Logger, service auth.Service) *Controller {
	return &Controller{logger: logger, service: service, validate: validator.New()}
}

func (ctrl *Controller) Register(c echo.Context) error {
	var req auth.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("name, email, dan password (min. 8 karakter) wajib diisi dengan benar"))
	}

	user, err := ctrl.service.Register(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success("registrasi berhasil", user))
}

func (ctrl *Controller) Login(c echo.Context) error {
	var req auth.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("email dan password wajib diisi"))
	}

	result, err := ctrl.service.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("login berhasil", result))
}
