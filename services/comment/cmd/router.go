package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/comment/handler"
	"github.com/Starostina-elena/investment_platform/services/comment/middleware"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /add/{proj_id}", middleware.AuthMiddleware(handler.CreateCommentHandler(h)))
	router.Handle("GET /read/{comm_id}", handler.GetCommentHandler(h))
	router.Handle("POST /edit/{comm_id}", middleware.AuthMiddleware(handler.UpdateCommentHandler(h)))
	router.Handle("DELETE /delete/{comm_id}", middleware.AuthMiddleware(handler.DeleteCommentHandler(h)))
	router.Handle("GET /read/all/{proj_id}", handler.GetProjectCommentsHandler(h))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
