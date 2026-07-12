package queue

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/fahmi0sd/waste-banks-and-donations/service/queue"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	service  queue.Service
	validate *validator.Validate
}

func NewController(logger *slog.Logger, service queue.Service) *Controller {
	return &Controller{logger: logger, service: service, validate: validator.New()}
}

type createQueueRequest struct {
	LocationID    int     `json:"location_id" validate:"required"`
	PreferredTime *string `json:"preferred_time,omitempty"`
}

func (ctrl *Controller) Create(c echo.Context) error {
	var req createQueueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("location_id wajib diisi"))
	}

	var preferredTime *time.Time
	if req.PreferredTime != nil && *req.PreferredTime != "" {
		t, err := time.Parse(time.RFC3339, *req.PreferredTime)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Error("format preferred_time harus RFC3339, contoh: 2026-07-12T09:00:00Z"))
		}
		preferredTime = &t
	}

	userID := auth.UserID(c)
	q, err := ctrl.service.CreateQueue(userID, req.LocationID, preferredTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success("antrian berhasil dibuat", q))
}

func (ctrl *Controller) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id antrian tidak valid"))
	}

	userID := auth.UserID(c)
	q, err := ctrl.service.GetQueue(id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("detail antrian", q))
}

func (ctrl *Controller) AdminList(c echo.Context) error {
	adminID := auth.UserID(c)
	list, err := ctrl.service.AdminListQueue(adminID)
	if err != nil {
		return c.JSON(http.StatusForbidden, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("daftar antrian menunggu", list))
}

func (ctrl *Controller) AdminVerify(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id antrian tidak valid"))
	}

	adminID := auth.UserID(c)
	q, err := ctrl.service.AdminVerifyQueue(adminID, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("antrian berhasil diverifikasi", q))
}
