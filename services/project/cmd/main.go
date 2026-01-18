package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/project/clients"
	"github.com/Starostina-elena/investment_platform/services/project/handler"
	"github.com/Starostina-elena/investment_platform/services/project/repo"
	"github.com/Starostina-elena/investment_platform/services/project/service"
	"github.com/Starostina-elena/investment_platform/services/project/storage"
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

func openOrgClient(log slog.Logger) *clients.OrgClient {
	return clients.NewOrgClient(os.Getenv("ORG_SERVICE_URL"), log)
}

func openMinioStorage(log slog.Logger) *storage.MinioStorage {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	pictureBucket := os.Getenv("MINIO_BUCKET")
	if pictureBucket == "" {
		pictureBucket = "projects"
	}

	minioStorage, err := storage.NewMinioStorage(endpoint, accessKey, secretKey, useSSL, pictureBucket)
	if err != nil {
		log.Error("failed to initialize minio storage", "error", err)
		os.Exit(1)
	}
	return minioStorage
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger.Info("starting project service")

	db := openDB()
	defer db.Close()

	repo := repo.NewRepo(db, *logger)

	orgClient := openOrgClient(*logger)
	minioStorage := openMinioStorage(*logger)
	service := service.NewService(repo, orgClient, minioStorage, *logger)

	handler := handler.NewHandler(service, *logger)

	router := getRouter(handler)

	appPort := os.Getenv("APP_PORT")
	srv := &http.Server{Addr: ":" + appPort, Handler: router}

	logger.Info("project service listening on :" + appPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down project service...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
