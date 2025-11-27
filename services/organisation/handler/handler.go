package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/organisation/service"
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler {
	return &Handler{service: s, log: log}
}

func CreateOrgHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		o, err := h.service.Create(r.Context(), req.Name)
		if err != nil {
			h.log.Error("failed to create org", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(o)
	}
}

func GetOrgHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.log.Error("invalid org id", "id", idStr, "error", err)
			http.Error(w, "invalid org id", http.StatusBadRequest)
			return
		}
		o, err := h.service.Get(r.Context(), id)
		if err != nil {
			h.log.Error("org not found", "id", id, "error", err)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(o)
	}
}
