package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/organisation/handler"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /orgs", handler.CreateOrgHandler(h))
	router.Handle("GET /orgs/{id}", handler.GetOrgHandler(h))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
