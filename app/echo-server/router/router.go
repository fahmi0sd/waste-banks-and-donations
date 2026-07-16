package router

import (
	"net/http"

	authCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/auth"
	backupCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/backup"
	calculatorCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/calculator"
	campaignCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/campaign"
	categoryCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/category"
	locationCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/location"
	notificationCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/notification"
	priceCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/price"
	queueCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/queue"
	reportCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/report"
	userCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/user"
	walletCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/wallet"
	wastetransactionCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/waste-transaction"
	"github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Controllers struct {
	Auth             *authCtrl.Controller
	User             *userCtrl.Controller
	Location         *locationCtrl.Controller
	Category         *categoryCtrl.Controller
	Price            *priceCtrl.Controller
	Report           *reportCtrl.Controller
	Queue            *queueCtrl.Controller
	Calculator       *calculatorCtrl.Controller
	WasteTransaction *wastetransactionCtrl.Controller
	Notification     *notificationCtrl.Controller
	Wallet           *walletCtrl.Controller
	Campaign         *campaignCtrl.Controller
	Backup           *backupCtrl.Controller
}

func RegisterPath(e *echo.Echo, jwtSecret string, db *gorm.DB, ctrls Controllers) {
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	jwtMiddleware := middleware.JWTMiddleware2(jwtSecret)

	// Auth
	e.POST("/auth/register", ctrls.Auth.Register)
	e.POST("/auth/login", ctrls.Auth.Login)

	// Profil
	e.GET("/users/me", ctrls.User.Me, jwtMiddleware)
	e.PATCH("/users/me", ctrls.User.UpdateMe, jwtMiddleware)

	// CRUD lokasi bank sampah
	e.GET("/locations", ctrls.Location.List)
	e.GET("/locations/:id", ctrls.Location.Get)
	e.POST("/locations", ctrls.Location.Create, jwtMiddleware)
	e.PATCH("/locations/:id", ctrls.Location.Update, jwtMiddleware)
	e.DELETE("/locations/:id", ctrls.Location.Delete, jwtMiddleware)
	// toggle buka/tutup lokasi
	e.PATCH("/locations/:id/toggle", ctrls.Location.Toggle, jwtMiddleware)

	// CRUD kategori sampah
	e.GET("/categories", ctrls.Category.List)
	e.GET("/categories/:id", ctrls.Category.Get)
	e.POST("/categories", ctrls.Category.Create, jwtMiddleware)
	e.PATCH("/categories/:id", ctrls.Category.Update, jwtMiddleware)
	e.DELETE("/categories/:id", ctrls.Category.Delete, jwtMiddleware)

	// Manajemen harga dengan histori
	e.GET("/categories/:id/price", ctrls.Price.ActivePrice)
	e.GET("/categories/:id/price/history", ctrls.Price.History, jwtMiddleware)
	e.POST("/categories/:id/price", ctrls.Price.SetPrice, jwtMiddleware)

	// Laporan transaksi (agregasi per lokasi/tanggal)
	e.GET("/admin/reports/transactions", ctrls.Report.TransactionReport, jwtMiddleware)

	// CRUD akun admin/user oleh master admin
	e.GET("/admin/users", ctrls.User.AdminList, jwtMiddleware)
	e.POST("/admin/users", ctrls.User.AdminCreate, jwtMiddleware)
	e.PATCH("/admin/users/:id", ctrls.User.AdminUpdate, jwtMiddleware)
	e.DELETE("/admin/users/:id", ctrls.User.AdminDelete, jwtMiddleware)

	// Calculator & Queue
	e.POST("/calculator/simulate", ctrls.Calculator.Simulate, jwtMiddleware)

	e.POST("/queue", ctrls.Queue.Create, jwtMiddleware)
	e.GET("/queue/:id", ctrls.Queue.Get, jwtMiddleware)
	e.GET("/admin/queue", ctrls.Queue.AdminList, jwtMiddleware)
	e.PATCH("/admin/queue/:id/verify", ctrls.Queue.AdminVerify, jwtMiddleware)

	// Waste Transaction
	e.POST("/admin/waste-transactions", ctrls.WasteTransaction.Create, jwtMiddleware)
	e.GET("/waste-transactions/:id", ctrls.WasteTransaction.Get, jwtMiddleware)
	e.GET("/users/me/waste-transactions", ctrls.WasteTransaction.MyHistory, jwtMiddleware)
	e.GET("/admin/locations/:id/waste-transactions", ctrls.WasteTransaction.LocationHistory, jwtMiddleware)
	e.PATCH("/admin/waste-transactions/:id/status", ctrls.WasteTransaction.UpdateStatus, jwtMiddleware)

	// Wallet & Withdraw
	e.GET("/users/me/wallet", ctrls.Wallet.GetWallet, jwtMiddleware)
	e.GET("/users/me/wallet/transactions", ctrls.Wallet.GetTransactions, jwtMiddleware)
	e.POST("/users/me/wallet/withdraw", ctrls.Wallet.Withdraw, jwtMiddleware)

	// Campaigns
	// Public
	e.GET("/campaigns", ctrls.Campaign.List)
	e.GET("/campaigns/:id", ctrls.Campaign.Get)
	e.GET("/campaigns/:id/updates", ctrls.Campaign.ListUpdate)
	// Master Admin
	e.POST("/master-admin/campaigns", ctrls.Campaign.Create, jwtMiddleware)
	e.PATCH("/master-admin/campaigns/:id", ctrls.Campaign.Update, jwtMiddleware)
	e.DELETE("/master-admin/campaigns/:id", ctrls.Campaign.Delete, jwtMiddleware)
	e.POST("/master-admin/campaigns/:id/updates", ctrls.Campaign.CreateUpdate, jwtMiddleware)
	// Donation
	e.POST("/campaigns/:id/donate", ctrls.Campaign.Donate, jwtMiddleware)
	e.GET("/users/me/donations", ctrls.Campaign.MyDonations, jwtMiddleware)

	// Notification
	e.GET("/users/me/notifications", ctrls.Notification.MyNotifications, jwtMiddleware)

	// Backup
	e.POST("/master-admin/backup/trigger", ctrls.Backup.Trigger, jwtMiddleware)
	e.GET("/master-admin/backup/history", ctrls.Backup.History, jwtMiddleware)

}
