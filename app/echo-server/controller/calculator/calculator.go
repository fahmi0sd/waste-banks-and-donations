package calculator

import (
	"log/slog"
	"net/http"

	"github.com/fahmi0sd/waste-banks-and-donations/service/calculator"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	service  calculator.Service
	validate *validator.Validate
}

func NewController(logger *slog.Logger, service calculator.Service) *Controller {
	return &Controller{logger: logger, service: service, validate: validator.New()}
}

func (ctrl *Controller) Simulate(c echo.Context) error {
	var req calculator.SimulateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("items wajib diisi minimal 1, tiap item wajib category_id dan weight > 0"))
	}

	result, err := ctrl.service.Simulate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("simulasi berhasil", result))
}
