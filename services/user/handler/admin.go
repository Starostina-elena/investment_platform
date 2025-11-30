package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Starostina-elena/investment_platform/services/user/middleware"
)

func SetAdminHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_requested := middleware.FromContext(r.Context())
		if !user_requested.Admin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		idStr := r.PathValue("user_id")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			h.log.Error("invalid user id", "id", idStr, "error", err)
			http.Error(w, "Некорретный id", http.StatusBadRequest)
			return
		}

		isAdminStr := strings.TrimSpace(r.URL.Query().Get("admin"))
		isAdmin, err := strconv.ParseBool(isAdminStr)
		if err != nil {
			h.log.Error("invalid admin value", "value", isAdminStr, "error", err)
			http.Error(w, "Некорретное значение admin", http.StatusBadRequest)
			return
		}

		err = h.service.SetAdmin(r.Context(), userId, isAdmin)
		if err != nil {
			h.log.Error("failed to set admin status", "user_id", userId, "is_admin", isAdmin, "error", err)
			http.Error(w, "Не удалось изменить статус администратора", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
