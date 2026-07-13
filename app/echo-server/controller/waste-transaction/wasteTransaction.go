package wastetransaction

import (
	"log/slog"
	"net/http"
	"strconv"

	wastetransaction "github.com/fahmi0sd/waste-banks-and-donations/service/waste-transaction"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	service  wastetransaction.Service
	validate *validator.Validate
}

func NewController(logger *slog.Logger, service wastetransaction.Service) *Controller {
	return &Controller{logger: logger, service: service, validate: validator.New()}
}

type createItemRequest struct {
	CategoryID int     `json:"category_id" validate:"required"`
	Weight     float64 `json:"weight" validate:"required,gt=0"`
}

type createRequest struct {
	QueueID int                 `json:"queue_id" validate:"required"`
	Items   []createItemRequest `json:"items" validate:"required,min=1,dive"`
}

func (ctrl *Controller) Create(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("queue_id wajib diisi, items minimal 1, tiap item wajib category_id dan weight > 0"))
	}

	items := make([]wastetransaction.ItemInput, len(req.Items))
	for i, it := range req.Items {
		items[i] = wastetransaction.ItemInput{CategoryID: it.CategoryID, Weight: it.Weight}
	}

	adminID := auth.UserID(c)
	tx, err := ctrl.service.Create(adminID, req.QueueID, items)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success("transaksi berhasil dibuat", tx))
}

func (ctrl *Controller) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id transaksi tidak valid"))
	}

	requesterID := auth.UserID(c)
	tx, err := ctrl.service.GetByID(id, requesterID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("detail transaksi", tx))
}

func (ctrl *Controller) MyHistory(c echo.Context) error {
	userID := auth.UserID(c)
	list, err := ctrl.service.MyHistory(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("riwayat transaksi", list))
}

func (ctrl *Controller) LocationHistory(c echo.Context) error {
	locationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id lokasi tidak valid"))
	}

	adminID := auth.UserID(c)
	list, err := ctrl.service.LocationHistory(adminID, locationID)
	if err != nil {
		return c.JSON(http.StatusForbidden, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("riwayat transaksi lokasi", list))
}

type updateStatusRequest struct {
	Status string `json:"status" validate:"required,eq=cancelled"`
}

func (ctrl *Controller) UpdateStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id transaksi tidak valid"))
	}

	var req updateStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("status hanya boleh diisi 'cancelled' — endpoint ini tidak untuk mengubah nominal/berat"))
	}

	adminID := auth.UserID(c)
	tx, err := ctrl.service.Cancel(adminID, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("transaksi berhasil dibatalkan", tx))
}
