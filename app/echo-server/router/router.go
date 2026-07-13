package router

import (
	"net/http"

	calculatorCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/calculator"
	queueCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/queue"
	wastetransactionCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/waste-transaction"
	"github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterPath(e *echo.Echo, jwtSecret string, db *gorm.DB, ctrlQueue *queueCtrl.Controller, ctrlCalculator *calculatorCtrl.Controller, ctrlWasteTransaction *wastetransactionCtrl.Controller) {
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	jwtMiddleware := middleware.JWTMiddleware2(jwtSecret)

	e.POST("/calculator/simulate", ctrlCalculator.Simulate, jwtMiddleware)

	e.POST("/queue", ctrlQueue.Create, jwtMiddleware)
	e.GET("/queue/:id", ctrlQueue.Get, jwtMiddleware)
	e.GET("/admin/queue", ctrlQueue.AdminList, jwtMiddleware)
	e.PATCH("/admin/queue/:id/verify", ctrlQueue.AdminVerify, jwtMiddleware)

	e.POST("/admin/waste-transactions", ctrlWasteTransaction.Create, jwtMiddleware)
	e.GET("/waste-transactions/:id", ctrlWasteTransaction.Get, jwtMiddleware)
	e.GET("/users/me/waste-transactions", ctrlWasteTransaction.MyHistory, jwtMiddleware)
	e.GET("/admin/locations/:id/waste-transactions", ctrlWasteTransaction.LocationHistory, jwtMiddleware)
	e.PATCH("/admin/waste-transactions/:id/status", ctrlWasteTransaction.UpdateStatus, jwtMiddleware)
}
