package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/user/middleware"
)

func GetActiveInvestmentsHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		investments, err := h.service.GetActiveInvestments(r.Context(), claims.UserID)
		if err != nil {
			h.log.Error("failed to get active investments", "user_id", claims.UserID, "error", err)
			http.Error(w, "Ошибка при получении активных инвестиций", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(investments)
	}
}

func GetArchivedInvestmentsHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		investments, err := h.service.GetArchivedInvestments(r.Context(), claims.UserID)
		if err != nil {
			h.log.Error("failed to get archived investments", "user_id", claims.UserID, "error", err)
			http.Error(w, "Ошибка при получении архивных инвестиций", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(investments)
	}
}
