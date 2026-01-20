package handler

import (
	"encoding/json"
	"net/http"
)

func ChangeBalanceHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Тут в идеале проверка на внутренний токен/секрет
		var req struct {
			UserID int     `json:"id"`
			Delta  float64 `json:"delta"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		err := h.service.ChangeBalance(r.Context(), req.UserID, req.Delta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
