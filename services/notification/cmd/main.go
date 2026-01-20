package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Starostina-elena/investment_platform/services/notification/handler"
	"github.com/Starostina-elena/investment_platform/services/notification/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger.Info("starting notification service")

	emailService := service.NewEmailService(*logger)
	h := handler.NewHandler(emailService, *logger)
	router := getRouter(h)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8083"
	}
	srv := &http.Server{Addr: ":" + appPort, Handler: router}
	logger.Info("notification service listening on :" + appPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down notification service...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
