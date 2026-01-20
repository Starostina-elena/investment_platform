package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Starostina-elena/investment_platform/services/payment/clients"
	"github.com/Starostina-elena/investment_platform/services/payment/handler"
	"github.com/Starostina-elena/investment_platform/services/payment/repo"
	"github.com/Starostina-elena/investment_platform/services/payment/service"
	"github.com/Starostina-elena/investment_platform/services/payment/yookassa"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// DB Connection
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

	r := repo.NewRepo(db)
	yc := yookassa.NewClient()
	tc := clients.NewTransactionClient()
	svc := service.NewService(r, yc, tc, *logger)
	h := handler.NewHandler(svc)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := svc.ProcessPendingPayments(ctx); err != nil {
				logger.Error("error processing pending payments", "error", err)
			}
			cancel()
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /pay/init", h.InitPaymentHandler)
	mux.HandleFunc("POST /pay/webhook", h.WebhookHandler)
	mux.HandleFunc("POST /pay/check", h.CheckPaymentHandler)

	logger.Info("payment service listening on :8106")
	http.ListenAndServe(":8106", mux)
}
