package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"

	"github.com/Starostina-elena/investment_platform/services/daemon/jobs"
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
	logger.Info("starting daemon service")

	db := openDB()
	defer db.Close()

	c := cron.New()

	expiredJob := jobs.NewExpiredProjectsJob(db, logger)
	_, err := c.AddFunc("0 0 * * *", expiredJob.Run)
	if err != nil {
		log.Fatalf("failed to add expired projects cron job: %v", err)
	}

	recalculateJob := jobs.NewRecalculatePaybackJob(db, logger)
	_, err = c.AddFunc("0 0 * * *", recalculateJob.Run)
	if err != nil {
		log.Fatalf("failed to add recalculate payback cron job: %v", err)
	}

	c.Start()
	logger.Info("daemon service started, cron jobs scheduled")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down daemon service...")
	c.Stop()
}
