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

func NewHandler(s service.Service, log slog.Logger) *Handler {
	return &Handler{service: s, log: log}
}

// InvestHandler обрабатывает запросы на инвестирование
func InvestHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ожидаем структуру, которую отправляет фронтенд
		var req struct {
			UserID    int     `json:"user_id"`
			ProjectID int     `json:"project_id"`
			Amount    float64 `json:"amount"`
			Method    string  `json:"method"` // "sbp" или "yookassa"
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode request", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Валидация
		if req.Amount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}

		// Вызываем новый метод сервиса Invest вместо старого Create
		tx, err := h.service.Invest(r.Context(), req.UserID, req.ProjectID, req.Amount, req.Method)
		if err != nil {
			h.log.Error("failed to invest", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(tx)
	}
}
