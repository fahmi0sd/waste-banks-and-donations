package notification

import (
	"log/slog"
	"net/http"

	"github.com/fahmi0sd/waste-banks-and-donations/service/notification"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger  *slog.Logger
	service notification.Service
}

func NewController(logger *slog.Logger, service notification.Service) *Controller {
	return &Controller{logger: logger, service: service}
}

func (ctrl *Controller) MyNotifications(c echo.Context) error {
	userID := auth.UserID(c)
	list, err := ctrl.service.MyNotifications(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("riwayat notifikasi", list))
}
