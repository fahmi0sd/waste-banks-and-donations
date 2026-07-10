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
	"github.com/fahmi0sd/waste-banks-and-donations/app/echo-server/router"
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

	router.RegisterPath(e, jwtSecret, database)

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
