package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/project/handler"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /projects", handler.CreateProjectHandler(h))
	router.Handle("GET /projects/{id}", handler.GetProjectHandler(h))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
