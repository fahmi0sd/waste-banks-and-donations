package report

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fahmi0sd/waste-banks-and-donations/service/report"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger  *slog.Logger
	service report.Service
}

func NewController(logger *slog.Logger, service report.Service) *Controller {
	return &Controller{logger: logger, service: service}
}

// TransactionReport menangani GET /admin/reports/transactions
// dengan query param opsional: location_id, date_from, date_to (format YYYY-MM-DD).
func (ctrl *Controller) TransactionReport(c echo.Context) error {
	requesterID := auth.UserID(c)

	var filter report.Filter
	if raw := c.QueryParam("location_id"); raw != "" {
		id, err := strconv.Atoi(raw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Error("location_id tidak valid"))
		}
		filter.LocationID = &id
	}
	if raw := c.QueryParam("date_from"); raw != "" {
		filter.DateFrom = &raw
	}
	if raw := c.QueryParam("date_to"); raw != "" {
		filter.DateTo = &raw
	}

	list, err := ctrl.service.TransactionReport(requesterID, filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("laporan transaksi", list))
}
