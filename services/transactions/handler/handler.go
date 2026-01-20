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

func WithdrawHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Withdrawal temporarily disabled. Use deposit to add funds.", http.StatusNotImplemented)
	}
}

// CreateDepositHandler создает платеж для пополнения баланса
func CreateDepositHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID    int     `json:"user_id"`
			Amount    float64 `json:"amount"`
			ReturnURL string  `json:"return_url"` // URL куда вернется пользователь после оплаты
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode deposit request", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Amount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}

		if req.ReturnURL == "" {
			req.ReturnURL = "http://localhost:3000/profile" // URL по умолчанию
		}

		paymentID, confirmationURL, err := h.service.CreateDeposit(r.Context(), req.UserID, req.Amount, req.ReturnURL)
		if err != nil {
			h.log.Error("failed to create deposit", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			PaymentID       string `json:"payment_id"`
			ConfirmationURL string `json:"confirmation_url"`
			Message         string `json:"message"`
		}{
			PaymentID:       paymentID,
			ConfirmationURL: confirmationURL,
			Message:         "Перейдите по ссылке для оплаты",
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}
}

// CheckDepositHandler проверяет статус платежа
func CheckDepositHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			PaymentID string `json:"payment_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode check request", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.PaymentID == "" {
			http.Error(w, "payment_id is required", http.StatusBadRequest)
			return
		}

		tx, err := h.service.CheckDeposit(r.Context(), req.PaymentID)
		if err != nil {
			h.log.Error("failed to check deposit", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(tx)
	}
}
