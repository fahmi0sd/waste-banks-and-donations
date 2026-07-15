package router

import (
	"net/http"

	authCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/auth"
	calculatorCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/calculator"
	categoryCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/category"
	locationCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/location"
	queueCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/queue"
	userCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/user"
	"github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Controllers struct {
	Auth       *authCtrl.Controller
	User       *userCtrl.Controller
	Location   *locationCtrl.Controller
	Category   *categoryCtrl.Controller
	Queue      *queueCtrl.Controller
	Calculator *calculatorCtrl.Controller
}

func RegisterPath(e *echo.Echo, jwtSecret string, db *gorm.DB, ctrls Controllers) {
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	jwtMiddleware := middleware.JWTMiddleware2(jwtSecret)

	// M2 #1, #2 - Auth (publik)
	e.POST("/auth/register", ctrls.Auth.Register)
	e.POST("/auth/login", ctrls.Auth.Login)

	// M2 #4 - Profil (butuh token, middleware JWT dari #3)
	e.GET("/users/me", ctrls.User.Me, jwtMiddleware)
	e.PATCH("/users/me", ctrls.User.UpdateMe, jwtMiddleware)

	// M3 #5 - CRUD lokasi bank sampah
	e.GET("/locations", ctrls.Location.List)
	e.GET("/locations/:id", ctrls.Location.Get)
	e.POST("/locations", ctrls.Location.Create, jwtMiddleware)
	e.PATCH("/locations/:id", ctrls.Location.Update, jwtMiddleware)
	e.DELETE("/locations/:id", ctrls.Location.Delete, jwtMiddleware)

	// M3 #6 - CRUD kategori sampah
	e.GET("/categories", ctrls.Category.List)
	e.GET("/categories/:id", ctrls.Category.Get)
	e.POST("/categories", ctrls.Category.Create, jwtMiddleware)
	e.PATCH("/categories/:id", ctrls.Category.Update, jwtMiddleware)
	e.DELETE("/categories/:id", ctrls.Category.Delete, jwtMiddleware)

	// Calculator & Queue (sudah ada sebelumnya)
	e.POST("/calculator/simulate", ctrls.Calculator.Simulate, jwtMiddleware)

	e.POST("/queue", ctrls.Queue.Create, jwtMiddleware)
	e.GET("/queue/:id", ctrls.Queue.Get, jwtMiddleware)
	e.GET("/admin/queue", ctrls.Queue.AdminList, jwtMiddleware)
	e.PATCH("/admin/queue/:id/verify", ctrls.Queue.AdminVerify, jwtMiddleware)
}
