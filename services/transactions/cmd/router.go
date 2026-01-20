package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/transactions/handler"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /transactions", handler.InvestHandler(h))

	// Пополнение баланса через ЮКассу
	router.Handle("POST /deposit", handler.CreateDepositHandler(h))      // Создать платеж
	router.Handle("POST /deposit/check", handler.CheckDepositHandler(h)) // Проверить статус

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	}))

	return router
}
