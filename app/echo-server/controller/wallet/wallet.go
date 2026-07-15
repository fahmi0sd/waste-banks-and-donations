package wallet

import (
	"log/slog"
	"net/http"

	walletService "github.com/fahmi0sd/waste-banks-and-donations/service/wallet"
	"github.com/fahmi0sd/waste-banks-and-donations/util/auth"
	"github.com/fahmi0sd/waste-banks-and-donations/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	service  walletService.Service
	validate *validator.Validate
}

func NewController(
	logger *slog.Logger,
	service walletService.Service,
) *Controller {

	return &Controller{
		logger:   logger,
		service:  service,
		validate: validator.New(),
	}
}

func (ctrl *Controller) GetWallet(c echo.Context) error {

	userID := auth.UserID(c)

	wallet, err := ctrl.service.GetWallet(userID)

	if err != nil {

		return c.JSON(
			http.StatusBadRequest,
			response.Error(err.Error()),
		)

	}

	return c.JSON(

		http.StatusOK,

		response.Success(
			"wallet berhasil diambil",
			wallet,
		),
	)
}

func (ctrl *Controller) GetTransactions(c echo.Context) error {

	userID := auth.UserID(c)

	transactions, err := ctrl.service.GetTransactions(userID)

	if err != nil {

		return c.JSON(
			http.StatusBadRequest,
			response.Error(err.Error()),
		)

	}

	return c.JSON(

		http.StatusOK,

		response.Success(
			"riwayat wallet berhasil diambil",
			transactions,
		),
	)
}

type withDrawRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

func (ctrl *Controller) Withdraw(c echo.Context) error {
	var req withDrawRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.Error("payload tidak valid"),
		)
	}

	if err := ctrl.validate.Struct(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.Error("amount wajib diisi dan lebih dari 0"),
		)
	}

	userID := auth.UserID(c)

	transaction, err := ctrl.service.Withdraw(userID, req.Amount)

	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.Error(err.Error()),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.Success("withdraw berhasil", transaction),
	)
}
