package handler

import (
	"encoding/json"
	"net/http"
)

// ChangeBalanceRequest - запрос на изменение баланса организации
type ChangeBalanceRequest struct {
	ID    int     `json:"id"`    // ID организации
	Delta float64 `json:"delta"` // Сумма изменения (может быть отрицательной)
}

// ChangeBalanceHandler обрабатывает изменение баланса организации
// POST /internal/balance
func ChangeBalanceHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeBalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.ID <= 0 {
			http.Error(w, "invalid org id", http.StatusBadRequest)
			return
		}

		// Вызываем сервис для обновления баланса
		if err := h.service.ChangeBalance(r.Context(), req.ID, req.Delta); err != nil {
			h.log.Error("failed to change balance", "org_id", req.ID, "delta", req.Delta, "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
