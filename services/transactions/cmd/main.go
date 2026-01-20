package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/transactions/clients"
	"github.com/Starostina-elena/investment_platform/services/transactions/handler"
	"github.com/Starostina-elena/investment_platform/services/transactions/repo"
	"github.com/Starostina-elena/investment_platform/services/transactions/service"
)

func openDB() *sqlx.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	return db
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("starting transactions service")

	db := openDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			// TODO maybe not ignore
		}
	}(db)

	// Исправляем имена переменных, чтобы не конфликтовали с пакетами
	repository := repo.NewRepo(db, *logger)
	projectClient := clients.NewProjectClient(*logger)
	svc := service.NewService(repository, projectClient, *logger)
	h := handler.NewHandler(svc, *logger)

	router := getRouter(h)

	appPort := os.Getenv("APP_PORT")
	srv := &http.Server{Addr: ":" + appPort, Handler: router}
	logger.Info("transactions service listening on :" + appPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down transactions service...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
