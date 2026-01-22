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
	EntityType string  `json:"entity_type"`
	EntityID   int     `json:"entity_id"`
	Amount     float64 `json:"amount"`
	ReturnURL  string  `json:"return_url"`
}

func (h *Handler) InitPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req InitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	entityType := req.EntityType
	entityID := req.EntityID
	if entityType == "" {
		entityType = "user"
	}
	if entityID == 0 {
		http.Error(w, "entity_id must be set", http.StatusBadRequest)
		return
	}

	url, err := h.service.InitPayment(r.Context(), entityType, entityID, req.Amount, req.ReturnURL)
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

	if req.Event == "payment.succeeded" || req.Event == "payment.failed" {
		if err := h.service.ProcessWebhook(r.Context(), req.Event, req.Object); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	if req.Event == "payout.succeeded" || req.Event == "payout.failed" {
		if err := h.service.ProcessWithdrawalWebhook(r.Context(), req.Event, req.Object); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

type CheckPaymentRequest struct {
	PaymentID string `json:"payment_id"`
}

func (h *Handler) CheckPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.PaymentID == "" {
		http.Error(w, "payment_id is required", http.StatusBadRequest)
		return
	}

	payment, err := h.service.CheckPayment(r.Context(), req.PaymentID)
	if err != nil {
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

type InitWithdrawalRequest struct {
	EntityType        string  `json:"entity_type"`
	EntityID          int     `json:"entity_id"`
	Amount            float64 `json:"amount"`
	PayoutDestination string  `json:"payout_destination"`
}

func (h *Handler) InitWithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	var req InitWithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	entityType := req.EntityType
	entityID := req.EntityID
	if entityType == "" {
		entityType = "user"
	}
	if entityID == 0 {
		http.Error(w, "entity_id must be set", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}

	if req.PayoutDestination == "" {
		http.Error(w, "payout_destination is required", http.StatusBadRequest)
		return
	}

	withdrawalID, err := h.service.InitWithdrawal(r.Context(), entityType, entityID, req.Amount, req.PayoutDestination)
	if err != nil {
		http.Error(w, "failed to init withdrawal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"withdrawal_id": withdrawalID})
}

type CheckWithdrawalRequest struct {
	WithdrawalID string `json:"withdrawal_id"`
}

func (h *Handler) CheckWithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckWithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.WithdrawalID == "" {
		http.Error(w, "withdrawal_id is required", http.StatusBadRequest)
		return
	}

	withdrawal, err := h.service.CheckWithdrawal(r.Context(), req.WithdrawalID)
	if err != nil {
		http.Error(w, "withdrawal not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawal)
}
