package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/comment/service"
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler {
	return &Handler{service: s, log: log}
}

func CreateCommentHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ProjectID int    `json:"project_id"`
			Body      string `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		c, err := h.service.Create(r.Context(), req.ProjectID, req.Body)
		if err != nil {
			h.log.Error("failed to create comment", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(c)
	}
}

func GetCommentHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.log.Error("invalid comment id", "id", idStr, "error", err)
			http.Error(w, "invalid comment id", http.StatusBadRequest)
			return
		}

		c, err := h.service.Get(r.Context(), id)
		if err != nil {
			h.log.Error("failed to get comment", "id", id, "error", err)
			http.Error(w, "No such comment", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(c)
	}
}
