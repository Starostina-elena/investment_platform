package handler

import (
	"encoding/json"
	"net/http"
)

type ChangeBalanceRequest struct {
	ID    int     `json:"id"`
	Delta float64 `json:"delta"`
}

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

		if err := h.service.ChangeBalance(r.Context(), req.ID, req.Delta); err != nil {
			h.log.Error("failed to change balance", "org_id", req.ID, "delta", req.Delta, "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
