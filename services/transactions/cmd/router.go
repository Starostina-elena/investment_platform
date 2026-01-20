package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/transactions/handler"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /transactions", handler.InvestHandler(h))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong")) // Игнорируем ошибку записи
	}))

	return router
}
