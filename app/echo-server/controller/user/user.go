package user

import (
	"log/slog"
	"net/http"
	"strconv"

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

// --- M3 #10: CRUD akun admin/user oleh master admin ---

func (ctrl *Controller) AdminList(c echo.Context) error {
	requesterID := auth.UserID(c)

	list, err := ctrl.service.ListAccounts(requesterID)
	if err != nil {
		return c.JSON(http.StatusForbidden, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("daftar akun", list))
}

func (ctrl *Controller) AdminCreate(c echo.Context) error {
	requesterID := auth.UserID(c)

	var req user.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("name, email, password (min. 8 karakter), dan role wajib diisi dengan benar"))
	}

	created, err := ctrl.service.CreateAccount(requesterID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("akun berhasil dibuat", created))
}

func (ctrl *Controller) AdminUpdate(c echo.Context) error {
	requesterID := auth.UserID(c)

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id akun tidak valid"))
	}

	var req user.UpdateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("payload tidak valid"))
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("data yang dikirim tidak valid"))
	}

	updated, err := ctrl.service.UpdateAccount(requesterID, targetID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("akun berhasil diperbarui", updated))
}

func (ctrl *Controller) AdminDelete(c echo.Context) error {
	requesterID := auth.UserID(c)

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id akun tidak valid"))
	}

	if err := ctrl.service.DeleteAccount(requesterID, targetID); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("akun berhasil dihapus", nil))
}
