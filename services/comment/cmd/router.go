package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/comment/handler"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /comments", handler.CreateCommentHandler(h))
	router.Handle("GET /comments/{id}", handler.GetCommentHandler(h))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
