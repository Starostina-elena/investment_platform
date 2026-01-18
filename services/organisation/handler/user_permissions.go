package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
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
