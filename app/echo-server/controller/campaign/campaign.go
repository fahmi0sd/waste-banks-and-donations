package campaign

import (
	"log/slog"
	"net/http"
	"strconv"

	campaignService "github.com/fahmi0sd/waste-banks-and-donations/service/campaign"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger  *slog.Logger
	service campaignService.Service
}

func NewController(logger *slog.Logger, service campaignService.Service) *Controller {
	return &Controller{
		logger:  logger,
		service: service,
	}
}

func (ctrl *Controller) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id tidak valid"))
	}
	result, err := ctrl.service.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("berhasil mendapatkan campaign", result))
}

func (ctrl *Controller) List(c echo.Context) error {
	result, err := ctrl.service.List()
	if err != nil {
		ctrl.logger.Error("failed get campaigns", "error", err)
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("berhasil mendapatkan daftar campaign", result))
}

func (ctrl *Controller) Create(c echo.Context) error {
	var req campaignService.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("request tidak valid"))
	}

	userID := auth.UserID(c)
	result, err := ctrl.service.Create(userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success("campaign berhasil dibuat", result))
}

func (ctrl *Controller) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id tidak valid"))
	}

	var req campaignService.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("request tidak valid"))
	}

	userID := auth.UserID(c)
	result, err := ctrl.service.Update(userID, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("campaign berhasil diperbarui", result))
}

func (ctrl *Controller) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id tidak valid"))
	}

	userID := auth.UserID(c)
	err = ctrl.service.Delete(userID, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("campaign berhasil ditutup", nil))
}

func (ctrl *Controller) CreateUpdate(c echo.Context) error {
	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id campaign tidak valid"))
	}

	var req campaignService.CreateUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("request tidak valid"))
	}

	userID := auth.UserID(c)
	result, err := ctrl.service.CreateUpdate(userID, campaignID, req)
	if err != nil {
		ctrl.logger.Error(
			"failed create campaign update",
			"campaign_id", campaignID,
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success("update campaign berhasil dibuat", result))
}

func (ctrl *Controller) ListUpdate(c echo.Context) error {
	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id campaign tidak valid"))
	}

	result, err := ctrl.service.ListUpdate(campaignID)
	if err != nil {
		ctrl.logger.Error(
			"failed get campaign updates",
			"campaign_id", campaignID,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("berhasil mendapatkan update campaign", result))
}

func (ctrl *Controller) Donate(c echo.Context) error {
	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("id campaign tidak valid"))
	}

	var req campaignService.DonateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("request tidak valid"))
	}

	userID := auth.UserID(c)
	err = ctrl.service.Donate(userID, campaignID, req)
	if err != nil {
		ctrl.logger.Error(
			"failed donate campaign",
			"user_id", userID,
			"campaign_id", campaignID,
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("donasi berhasil", nil))
}

func (ctrl *Controller) MyDonations(c echo.Context) error {
	userID := auth.UserID(c)
	result, err := ctrl.service.MyDonations(userID)
	if err != nil {
		ctrl.logger.Error(
			"failed get donation history",
			"user_id", userID,
			"error", err)
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("berhasil mendapatkan riwayat donasi", result))
}
