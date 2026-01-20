package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/payment/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

type InitRequest struct {
	UserID    int     `json:"user_id"`
	Amount    float64 `json:"amount"`
	ReturnURL string  `json:"return_url"`
}

func (h *Handler) InitPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req InitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	url, err := h.service.InitPayment(r.Context(), req.UserID, req.Amount, req.ReturnURL)
	if err != nil {
		http.Error(w, "failed to init payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"confirmation_url": url})
}

type WebhookRequest struct {
	Type   string                 `json:"type"`
	Event  string                 `json:"event"`
	Object map[string]interface{} `json:"object"`
}

func (h *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var req WebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Обрабатываем payment.succeeded
	if err := h.service.ProcessWebhook(r.Context(), req.Event, req.Object); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
