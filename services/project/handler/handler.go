package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/project/service"
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler { return &Handler{service: s, log: log} }

func CreateProjectHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		p, err := h.service.Create(r.Context(), req.Name)
		if err != nil {
			h.log.Error("failed to create project", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(p)
	}
}

func GetProjectHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.log.Error("invalid project id", "id", idStr, "error", err)
			http.Error(w, "invalid project id", http.StatusBadRequest)
		}

		p, err := h.service.Get(r.Context(), id)
		if err != nil {
			h.log.Error("project not found", "id", id, "error", err)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(p)
	}
}
