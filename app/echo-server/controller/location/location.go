package location

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fahmi0sd/waste-banks-and-donations/service/location"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	service  location.Service
	validate *validator.Validate
}

func NewController(logger *slog.Logger, service location.Service) *Controller {
	return &Controller{logger: logger, service: service, validate: validator.New()}
}

func (ctrl *Controller) Create(c echo.Context) error {
	requesterID := auth.UserID(c)

	var req location.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("name wajib diisi"))
	}

	created, err := ctrl.service.Create(requesterID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("lokasi berhasil dibuat", created))
}

func (ctrl *Controller) List(c echo.Context) error {
	list, err := ctrl.service.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("daftar lokasi", list))
}

func (ctrl *Controller) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id lokasi tidak valid"))
	}

	l, err := ctrl.service.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("detail lokasi", l))
}

func (ctrl *Controller) Update(c echo.Context) error {
	requesterID := auth.UserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id lokasi tidak valid"))
	}

	var req location.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("data yang dikirim tidak valid"))
	}

	updated, err := ctrl.service.Update(requesterID, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("lokasi berhasil diperbarui", updated))
}

func (ctrl *Controller) Delete(c echo.Context) error {
	requesterID := auth.UserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id lokasi tidak valid"))
	}

	if err := ctrl.service.Delete(requesterID, id); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("lokasi berhasil dihapus", nil))
}

// Toggle menangani PATCH /locations/:id/toggle — buka/tutup lokasi (M3 #9).
func (ctrl *Controller) Toggle(c echo.Context) error {
	requesterID := auth.UserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id lokasi tidak valid"))
	}

	updated, err := ctrl.service.ToggleOpen(requesterID, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("status lokasi berhasil diubah", updated))
}
