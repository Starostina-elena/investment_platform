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

	"github.com/Starostina-elena/investment_platform/services/user/handler"
	"github.com/Starostina-elena/investment_platform/services/user/repo"
	"github.com/Starostina-elena/investment_platform/services/user/service"
	"github.com/Starostina-elena/investment_platform/services/user/storage"
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

func openMinio() *storage.MinioStorage {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	bucketName := os.Getenv("MINIO_BUCKET")

	var minioStorage *storage.MinioStorage
	var err error

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		minioStorage, err = storage.NewMinioStorage(endpoint, accessKey, secretKey, useSSL, bucketName)
		if err == nil {
			log.Printf("Successfully connected to MinIO on attempt %d", i+1)
			return minioStorage
		}
		log.Printf("Failed to connect to MinIO (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("failed to connect to MinIO after %d attempts: %v", maxRetries, err)
	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("starting user service")

	db := openDB()
	defer db.Close()

	minioStorage := openMinio()

	repo := repo.NewRepo(db, *logger)
	service := service.NewService(repo, *minioStorage, *logger)
	handler := handler.NewHandler(service, *logger)

	router := getRouter(handler)

	appPort := os.Getenv("APP_PORT")
	srv := &http.Server{Addr: ":" + appPort, Handler: router}
	logger.Info("user service listening on :" + appPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down user service...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
