package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/user/handler"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /create", handler.CreateUserHandler(h))
	router.Handle("GET /{id}", handler.GetUserHandler(h))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
