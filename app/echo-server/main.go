package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmi0sd/go-utils/postgres"
	calculatorCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/calculator"
	notificationCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/notification"
	queueCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/queue"
	wastetransactionCtrl "github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/controller/waste-transaction"
	"github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/router"
	"github.com/fahmi0sd/waste-banks-and-donations/pkg"
	calculatorRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/calculator"
	notificationRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/notification"
	queueRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/queue"
	wastetransactionRepo "github.com/fahmi0sd/waste-banks-and-donations/repository/waste-transaction"
	calculatorSvc "github.com/fahmi0sd/waste-banks-and-donations/service/calculator"
	notificationSvc "github.com/fahmi0sd/waste-banks-and-donations/service/notification"
	queueSvc "github.com/fahmi0sd/waste-banks-and-donations/service/queue"
	wastetransactionSvc "github.com/fahmi0sd/waste-banks-and-donations/service/waste-transaction"
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

	router.RegisterPath(e, jwtSecret, database, qCtrl, calcCtrl, wtCtrl, notifCtrl)

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
