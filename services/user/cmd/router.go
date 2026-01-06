package main

import (
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/user/handler"
	"github.com/Starostina-elena/investment_platform/services/user/middleware"
)

func getRouter(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /create", handler.CreateUserHandler(h))
	router.Handle("POST /update", middleware.AuthMiddleware(handler.UpdateUserHandler(h)))
	router.Handle("GET /{id}", handler.GetUserHandler(h))

	router.Handle("POST /login", handler.LoginHandler(h))
	router.Handle("POST /refresh", handler.RefreshHandler(h))
	router.Handle("POST /logout", handler.LogoutHandler(h))

	router.Handle("POST /{user_id}/admin", middleware.AuthMiddleware(handler.SetAdminHandler(h)))
	router.Handle("POST /{user_id}/active", middleware.AuthMiddleware(handler.BanUserHandler(h)))

	router.Handle("POST /avatar/upload", middleware.AuthMiddleware(handler.UploadAvatarHandler(h)))
	router.Handle("DELETE /avatar/", middleware.AuthMiddleware(handler.DeleteAvatarHandler(h)))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
