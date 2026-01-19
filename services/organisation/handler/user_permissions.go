package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
	"github.com/Starostina-elena/investment_platform/services/organisation/middleware"
)

func CheckUserOrgPermissionHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgIDStr := r.PathValue("org_id")
		orgID, err := strconv.Atoi(orgIDStr)
		if err != nil {
			h.log.Error("invalid org id", "id", orgIDStr, "error", err)
			http.Error(w, "Incorrect org_id", http.StatusBadRequest)
			return
		}

		userIDStr := r.PathValue("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			h.log.Error("invalid user id", "id", userIDStr, "error", err)
			http.Error(w, "Incorrect user_id", http.StatusBadRequest)
			return
		}

		permission := r.PathValue("permission")

		allowed, err := h.service.CheckUserOrgPermission(r.Context(), orgID, userID, core.OrgPermission(permission))
		if err != nil {
			h.log.Error("failed to check user organisation permission", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		resp := map[string]bool{"allowed": allowed}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

type AddEmployeeReq struct {
	UserID     int  `json:"user_id"`
	OrgAccMgmt bool `json:"org_account_management"`
	MoneyMgmt  bool `json:"money_management"`
	ProjMgmt   bool `json:"project_management"`
}

func AddEmployeeHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		orgIDStr := r.PathValue("org_id")
		orgID, err := strconv.Atoi(orgIDStr)
		if err != nil {
			h.log.Error("invalid org id", "id", orgIDStr, "error", err)
			http.Error(w, "Incorrect org_id", http.StatusBadRequest)
			return
		}

		var req AddEmployeeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.UserID <= 0 {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		if !req.OrgAccMgmt && !req.MoneyMgmt && !req.ProjMgmt {
			http.Error(w, "At least one permission must be granted", http.StatusBadRequest)
			return
		}

		if err := h.service.AddEmployee(r.Context(), orgID, claims.UserID, req.UserID, req.OrgAccMgmt, req.MoneyMgmt, req.ProjMgmt); err != nil {
			if err == core.ErrNotAuthorized {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			h.log.Error("failed to add employee", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetOrgEmployeesHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgIDStr := r.PathValue("org_id")
		orgID, err := strconv.Atoi(orgIDStr)
		if err != nil {
			h.log.Error("invalid org id", "id", orgIDStr, "error", err)
			http.Error(w, "Incorrect org_id", http.StatusBadRequest)
			return
		}

		employees, err := h.service.GetOrgEmployees(r.Context(), orgID)
		if err != nil {
			h.log.Error("failed to get organisation employees", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if employees == nil {
			employees = []core.OrgEmployee{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(employees); err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}
