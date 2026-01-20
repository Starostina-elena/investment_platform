package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Starostina-elena/investment_platform/services/notification/core"
	"github.com/Starostina-elena/investment_platform/services/notification/service"
)

type Handler struct {
	emailService *service.EmailService
	log          slog.Logger
}

func NewHandler(emailService *service.EmailService, log slog.Logger) *Handler {
	return &Handler{emailService: emailService, log: log}
}

func SendEmailHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req core.EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode email request", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		h.log.Info("received email request", "email", req.Email, "type", req.Type, "project_name", req.ProjectName, "amount", req.Amount)

		if req.Email == "" || !strings.Contains(req.Email, "@") {
			h.log.Error("invalid email", "email", req.Email)
			http.Error(w, "invalid email", http.StatusBadRequest)
			return
		}
		if req.Type != core.NotifTypeDividends && req.Type != core.NotifTypeProjectClosed {
			h.log.Error("unknown notification type", "type", req.Type)
			http.Error(w, "unknown notification type", http.StatusBadRequest)
			return
		}
		if req.ProjectName == "" {
			h.log.Error("invalid project name")
			http.Error(w, "invalid project name", http.StatusBadRequest)
			return
		}
		if req.Type != core.NotifTypeProjectClosed && req.Amount <= 0 {
			h.log.Error("invalid amount", "amount", req.Amount, "type", req.Type)
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}

		if err := h.emailService.SendNotification(&req); err != nil {
			http.Error(w, "failed to send email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"sent"}`))
	}
}
