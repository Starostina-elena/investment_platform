package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/transactions/service"
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler { return &Handler{service: s, log: log} }

func CreateTransactionHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			From, To int
			Amount   float64
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		tx, err := h.service.Create(r.Context(), req.From, req.To, req.Amount)
		if err != nil {
			h.log.Error("failed to create tx", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(tx)
	}
}
