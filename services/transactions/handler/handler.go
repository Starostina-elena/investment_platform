package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/transactions/clients"
	"github.com/Starostina-elena/investment_platform/services/transactions/service"
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler {
	return &Handler{service: s, log: log}
}

// TransferHandler обрабатывает внутренние переводы
func TransferHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			FromType string  `json:"from_type"` // "user", "org", "project"
			FromID   int     `json:"from_id"`
			ToType   string  `json:"to_type"`
			ToID     int     `json:"to_id"`
			Amount   float64 `json:"amount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		tx, err := h.service.Transfer(
			r.Context(),
			clients.EntityType(req.FromType),
			clients.EntityType(req.ToType),
			req.FromID,
			req.ToID,
			req.Amount,
		)

		if err != nil {
			h.log.Error("transfer error", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError) // Или 400 в зависимости от типа ошибки
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(tx)
	}
}
