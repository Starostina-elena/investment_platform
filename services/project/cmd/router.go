package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/project/handler"
	"github.com/Starostina-elena/investment_platform/services/project/middleware"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /create", middleware.AuthMiddleware(handler.CreateProjectHandler(h)))
	router.Handle("GET /{id}", handler.GetProjectHandler(h))
	router.Handle("POST /{id}/update", middleware.AuthMiddleware(handler.UpdateProjectHandler(h)))

	router.Handle("GET /projects", handler.GetProjectListHandler(h))
	router.Handle("GET /projects/org/{creator_id}", handler.GetProjectsByCreatorHandler(h))

	router.Handle("POST /{id}/ban", middleware.AuthMiddleware(handler.BanProjectHandler(h)))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
