package price

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fahmi0sd/waste-banks-and-donations/service/price"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	service  price.Service
	validate *validator.Validate
}

func NewController(logger *slog.Logger, service price.Service) *Controller {
	return &Controller{logger: logger, service: service, validate: validator.New()}
}

// SetPrice menangani POST /categories/:id/price — menutup harga aktif lama
// dan menyimpan harga baru (M3 #7: manajemen harga dengan histori).
func (ctrl *Controller) SetPrice(c echo.Context) error {
	requesterID := auth.UserID(c)

	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id kategori tidak valid"))
	}

	var req price.SetPriceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("price_per_unit wajib diisi dan harus > 0"))
	}

	created, err := ctrl.service.SetPrice(requesterID, categoryID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("harga baru berhasil disimpan", created))
}

func (ctrl *Controller) ActivePrice(c echo.Context) error {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id kategori tidak valid"))
	}

	p, err := ctrl.service.ActivePrice(categoryID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("harga aktif", p))
}

func (ctrl *Controller) History(c echo.Context) error {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id kategori tidak valid"))
	}

	list, err := ctrl.service.History(categoryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("histori harga", list))
}
