package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/fahmi0sd/go-utils/postgres"
	authCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/auth"
	calculatorCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/calculator"
	categoryCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/category"
	locationCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/location"
	priceCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/price"
	queueCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/queue"
	reportCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/report"
	userCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/user"
	"github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/router"
	authRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/auth"
	calculatorRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/calculator"
	categoryRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/category"
	locationRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/location"
	priceRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/price"
	queueRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/queue"
	reportRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/report"
	userRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/user"
	authSvc "github.com/fahmi0sd/waste-banks-and-donations/service/auth"
	calculatorSvc "github.com/fahmi0sd/waste-banks-and-donations/service/calculator"
	categorySvc "github.com/fahmi0sd/waste-banks-and-donations/service/category"
	locationSvc "github.com/fahmi0sd/waste-banks-and-donations/service/location"
	priceSvc "github.com/fahmi0sd/waste-banks-and-donations/service/price"
	queueSvc "github.com/fahmi0sd/waste-banks-and-donations/service/queue"
	reportSvc "github.com/fahmi0sd/waste-banks-and-donations/service/report"
	userSvc "github.com/fahmi0sd/waste-banks-and-donations/service/user"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	// Database Connection
	database := postgres.GetPostgresConnection()
	logger.Info("postgres connected")

	// Config
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtTTLHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || jwtTTLHours <= 0 {
		jwtTTLHours = 24
	}

	// Auth (M2 #1, #2)
	aRepo := authRepo.NewGormRepository(database)
	aSvc := authSvc.NewService(logger, aRepo, jwtSecret, jwtTTLHours)
	aCtrl := authCtrl.NewController(logger, aSvc)

	// User / profil (M2 #4)
	uRepo := userRepo.NewGormRepository(database)
	uSvc := userSvc.NewService(logger, uRepo)
	uCtrl := userCtrl.NewController(logger, uSvc)

	// Location (M3 #5)
	lRepo := locationRepo.NewGormRepository(database)
	lSvc := locationSvc.NewService(logger, lRepo)
	lCtrl := locationCtrl.NewController(logger, lSvc)

	// Category (M3 #6)
	catRepo := categoryRepo.NewGormRepository(database)
	catSvc := categorySvc.NewService(logger, catRepo)
	catCtrl := categoryCtrl.NewController(logger, catSvc)

	// Price (M3 #7)
	priceRepository := priceRepo.NewGormRepository(database)
	priceService := priceSvc.NewService(logger, priceRepository)
	priceController := priceCtrl.NewController(logger, priceService)

	// Report (M11)
	repRepo := reportRepo.NewGormRepository(database)
	repSvc := reportSvc.NewService(logger, repRepo)
	repCtrl := reportCtrl.NewController(logger, repSvc)

	// Queue
	qRepo := queueRepo.NewGormRepository(database)
	qSvc := queueSvc.NewService(logger, qRepo)
	qCtrl := queueCtrl.NewController(logger, qSvc)

	// Calculator
	calcRepo := calculatorRepo.NewGormRepository(database)
	calcSvc := calculatorSvc.NewService(logger, calcRepo)
	calcCtrl := calculatorCtrl.NewController(logger, calcSvc)

	// Echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","level":"INFO","method":"${method}","uri":"${uri}","status":${status},"latency_human":"${latency_human}"}` + "\n",
	}))
	e.Pre(middleware.RemoveTrailingSlash())

	router.RegisterPath(e, jwtSecret, database, router.Controllers{
		Auth:       aCtrl,
		User:       uCtrl,
		Location:   lCtrl,
		Category:   catCtrl,
		Price:      priceController,
		Report:     repCtrl,
		Queue:      qCtrl,
		Calculator: calcCtrl,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	go func() {
		if err := e.Start(addr); err != http.ErrServerClosed {
			log.Fatal("server error: " + err.Error())
		}
	}()

	logger.Info("Bank Sampah & Donasi API running", slog.String("addr", addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("failed to shutdown server:", err)
	}

	logger.Info("server shutdown gracefully")
}
