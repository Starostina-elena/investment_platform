package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	dbpkg "investment_platform/db"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger.Info("starting app")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	logger.Info("Try to reach database on", "path", connStr)

	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		logger.Error("failed to connect to db", "error", err)
		return
	}
	defer db.Close()
	logger.Info("Connected to database")

	appDB := &dbpkg.DB{
		Log:  logger,
		Conn: db,
	}

	if err := appDB.Migrate(); err != nil {
		logger.Error("migration failed", "error", err)
		return
	}
	logger.Info("Database migrated successfully")

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
