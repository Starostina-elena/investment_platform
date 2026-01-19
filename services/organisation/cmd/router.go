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
	router.Handle("GET /{id}/full", middleware.AuthMiddleware(handler.GetFullOrgHandler(h)))
	router.Handle("POST /{org_id}/update", middleware.AuthMiddleware(handler.UpdateOrgHandler(h)))

	router.Handle("POST /{org_id}/avatar/upload", middleware.AuthMiddleware(handler.UploadAvatarHandler(h)))
	router.Handle("DELETE /{org_id}/avatar/", middleware.AuthMiddleware(handler.DeleteAvatarHandler(h)))

	router.Handle("POST /{org_id}/docs/{doc_type}", middleware.AuthMiddleware(handler.UploadDocHandler(h)))
	router.Handle("GET /{org_id}/docs/{doc_type}", middleware.AuthMiddleware(handler.DownloadDocHandler(h)))
	router.Handle("DELETE /{org_id}/docs/{doc_type}", middleware.AuthMiddleware(handler.DeleteDocHandler(h)))

	router.Handle("POST /{org_id}/active", middleware.AuthMiddleware(handler.BanOrgHandler(h)))

	router.Handle("GET /my", middleware.AuthMiddleware(handler.GetUserOrgsHandler(h)))

	router.Handle("GET /{org_id}/rights/{user_id}/{permission}", handler.CheckUserOrgPermissionHandler(h))
	router.Handle("POST /{org_id}/employees/add", middleware.AuthMiddleware(handler.AddEmployeeHandler(h)))
	router.Handle("GET /{org_id}/employees", handler.GetOrgEmployeesHandler(h))
	router.Handle("POST /{org_id}/employees/update", middleware.AuthMiddleware(handler.UpdateEmployeePermissionsHandler(h)))
	router.Handle("DELETE /{org_id}/employees/{user_id}/delete", middleware.AuthMiddleware(handler.DeleteEmployeeHandler(h)))
	router.Handle("POST /{org_id}/ownership/transfer/{new_owner_user_id}", middleware.AuthMiddleware(handler.TransferOwnershipHandler(h)))

	router.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	return router
}
