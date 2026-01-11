package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/organisation/handler"
	"github.com/Starostina-elena/investment_platform/services/organisation/middleware"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /create", middleware.AuthMiddleware(handler.CreateOrgHandler(h)))
	router.Handle("GET /{id}", handler.GetOrgHandler(h))
	router.Handle("POST /{org_id}/update", middleware.AuthMiddleware(handler.UpdateOrgHandler(h)))

	router.Handle("POST /{org_id}/avatar/upload", middleware.AuthMiddleware(handler.UploadAvatarHandler(h)))
	router.Handle("DELETE /{org_id}/avatar/", middleware.AuthMiddleware(handler.DeleteAvatarHandler(h)))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
