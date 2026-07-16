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
	"github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/router"
	"github.com/fahmi0sd/waste-banks-and-donations/pkg"
	authRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/auth"
	backupRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/backup"
	calculatorRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/calculator"
	categoryRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/category"
	locationRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/location"
	notificationRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/notification"
	priceRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/price"
	queueRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/queue"
	reportRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/report"
	userRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/user"
	walletRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/wallet"
	wastetransactionRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/waste-transaction"
	"github.com/fahmi0sd/waste-banks-and-donations/scheduler"
	authSvc "github.com/fahmi0sd/waste-banks-and-donations/service/auth"
	backupSvc "github.com/fahmi0sd/waste-banks-and-donations/service/backup"
	calculatorSvc "github.com/fahmi0sd/waste-banks-and-donations/service/calculator"
	categorySvc "github.com/fahmi0sd/waste-banks-and-donations/service/category"
	locationSvc "github.com/fahmi0sd/waste-banks-and-donations/service/location"
	notificationSvc "github.com/fahmi0sd/waste-banks-and-donations/service/notification"
	priceSvc "github.com/fahmi0sd/waste-banks-and-donations/service/price"
	queueSvc "github.com/fahmi0sd/waste-banks-and-donations/service/queue"
	reportSvc "github.com/fahmi0sd/waste-banks-and-donations/service/report"
	userSvc "github.com/fahmi0sd/waste-banks-and-donations/service/user"
	walletSvc "github.com/fahmi0sd/waste-banks-and-donations/service/wallet"
	wastetransactionSvc "github.com/fahmi0sd/waste-banks-and-donations/service/waste-transaction"

	campaignRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/campaign"

	campaignSvc "github.com/fahmi0sd/waste-banks-and-donations/service/campaign"
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
	mailjetAPIKey := os.Getenv("MAILJET_API_KEY")
	mailjetSecretKey := os.Getenv("MAILJET_SECRET_KEY")
	mailjetSenderEmail := os.Getenv("MAILJET_SENDER_EMAIL")
	mailjetSenderName := os.Getenv("MAILJET_SENDER_NAME")
	if mailjetSenderName == "" {
		mailjetSenderName = "Bank Sampah & Donasi"
	}
	mailjetBaseURL := os.Getenv("MAILJET_BASE_URL")
	if mailjetBaseURL == "" {
		mailjetBaseURL = "https://api.mailjet.com/v3.1/send"
	}

	// Auth
	aRepo := authRepo.NewGormRepository(database)
	aSvc := authSvc.NewService(logger, aRepo, jwtSecret, jwtTTLHours)
	aCtrl := authCtrl.NewController(logger, aSvc)

	// User / profil + admin user management
	uRepo := userRepo.NewGormRepository(database)
	uSvc := userSvc.NewService(logger, uRepo)
	uCtrl := userCtrl.NewController(logger, uSvc)

	// Location
	lRepo := locationRepo.NewGormRepository(database)
	lSvc := locationSvc.NewService(logger, lRepo)
	lCtrl := locationCtrl.NewController(logger, lSvc)

	// Category
	catRepo := categoryRepo.NewGormRepository(database)
	catSvc := categorySvc.NewService(logger, catRepo)
	catCtrl := categoryCtrl.NewController(logger, catSvc)

	// Price
	priceRepository := priceRepo.NewGormRepository(database)
	priceService := priceSvc.NewService(logger, priceRepository)
	priceController := priceCtrl.NewController(logger, priceService)

	// Report
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

	// Notification
	mailjetClient := pkg.NewMailjetClient(mailjetAPIKey, mailjetSecretKey, mailjetSenderEmail, mailjetSenderName, mailjetBaseURL)
	notifRepo := notificationRepo.NewGormRepository(database)
	notifSvc := notificationSvc.NewService(logger, notifRepo, mailjetClient)
	notifCtrl := notificationCtrl.NewController(logger, notifSvc)

	// Waste Transaction
	wtRepo := wastetransactionRepo.NewGormRepository(database)
	wtSvc := wastetransactionSvc.NewService(logger, wtRepo, database, notifSvc)
	wtCtrl := wastetransactionCtrl.NewController(logger, wtSvc)

	// Wallet
	wRepo := walletRepo.NewGormRepository(database)
	wSvc := walletSvc.NewService(logger, wRepo)
	wCtrl := walletCtrl.NewController(logger, wSvc)

	// Campaign
	campRepo := campaignRepo.NewGormRepository(database)
	campSvc := campaignSvc.NewService(logger, campRepo, database, notifSvc)
	campCtrl := campaignCtrl.NewController(logger, campSvc)

	// Backup
	bRepo := backupRepo.NewGormRepository(database)
	bSvc := backupSvc.NewService(logger, bRepo, database)
	bCtrl := backupCtrl.NewController(logger, bSvc)

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
		Auth:             aCtrl,
		User:             uCtrl,
		Location:         lCtrl,
		Category:         catCtrl,
		Price:            priceController,
		Report:           repCtrl,
		Queue:            qCtrl,
		Calculator:       calcCtrl,
		WasteTransaction: wtCtrl,
		Notification:     notifCtrl,
		Wallet:           wCtrl,
		Campaign:         campCtrl,
		Backup:           bCtrl,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	scheduler.StartBackupScheduler(logger, bSvc)

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
