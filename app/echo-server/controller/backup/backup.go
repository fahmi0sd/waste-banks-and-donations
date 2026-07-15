package backup

import (
	"log/slog"
	"net/http"

	backupService "github.com/fahmi0sd/waste-banks-and-donations/service/backup"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger  *slog.Logger
	service backupService.Service
}

func NewController(
	logger *slog.Logger,
	service backupService.Service,
) *Controller {

	return &Controller{
		logger:  logger,
		service: service,
	}
}

func (ctrl *Controller) Trigger(c echo.Context) error {

	userID := auth.UserID(c)

	result, err := ctrl.service.Trigger(userID)

	if err != nil {

		ctrl.logger.Error(
			"failed trigger backup",
			"user_id", userID,
			"error", err,
		)

		return c.JSON(
			http.StatusBadRequest,
			response.Error(err.Error()),
		)

	}

	return c.JSON(
		http.StatusOK,
		response.Success(
			"backup berhasil dibuat",
			result,
		),
	)
}

func (ctrl *Controller) History(c echo.Context) error {

	result, err := ctrl.service.History()

	if err != nil {

		ctrl.logger.Error(
			"failed get backup history",
			"error", err,
		)

		return c.JSON(
			http.StatusInternalServerError,
			response.Error(err.Error()),
		)

	}

	return c.JSON(
		http.StatusOK,
		response.Success(
			"history backup berhasil diambil",
			result,
		),
	)
}
