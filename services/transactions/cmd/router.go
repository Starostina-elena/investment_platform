package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/transactions/handler"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /transactions", handler.CreateTransactionHandler(h))
	// router.Handle("GET /transactions/{id}", handler.GetTransactionHandler(h))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
