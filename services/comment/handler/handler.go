package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Starostina-elena/investment_platform/services/comment/core"
	"github.com/Starostina-elena/investment_platform/services/comment/middleware"
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
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if claims.Banned {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		projIDStr := r.PathValue("proj_id")
		projID, err := strconv.Atoi(projIDStr)
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		var req struct {
			Body string `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		req.Body = strings.TrimSpace(req.Body)
		if len(req.Body) == 0 || len(req.Body) > 4096 {
			http.Error(w, "invalid comment body", http.StatusBadRequest)
			return
		}

		c, err := h.service.Create(r.Context(), projID, req.Body, claims.UserID)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(c)
	}
}

func GetCommentHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commIDStr := r.PathValue("comm_id")
		commID, err := strconv.Atoi(commIDStr)
		if err != nil {
			http.Error(w, "invalid comment id", http.StatusBadRequest)
			return
		}

		c, err := h.service.Get(r.Context(), commID)
		if err != nil {
			http.Error(w, "comment not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(c)
	}
}

func UpdateCommentHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if claims.Banned {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		commIDStr := r.PathValue("comm_id")
		commID, err := strconv.Atoi(commIDStr)
		if err != nil {
			http.Error(w, "invalid comment id", http.StatusBadRequest)
			return
		}

		comment, err := h.service.Get(r.Context(), commID)
		if err != nil {
			if err == core.ErrCommentNotFound {
				http.Error(w, "comment not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if comment.UserID != claims.UserID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		var req struct {
			Body string `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		req.Body = strings.TrimSpace(req.Body)
		if len(req.Body) == 0 || len(req.Body) > 4096 {
			http.Error(w, "invalid comment body", http.StatusBadRequest)
			return
		}

		c, err := h.service.Update(r.Context(), commID, req.Body)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(c)
	}
}

func DeleteCommentHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		commIDStr := r.PathValue("comm_id")
		commID, err := strconv.Atoi(commIDStr)
		if err != nil {
			http.Error(w, "invalid comment id", http.StatusBadRequest)
			return
		}

		comment, err := h.service.Get(r.Context(), commID)
		if err != nil {
			if err == core.ErrCommentNotFound {
				http.Error(w, "comment not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if comment.UserID != claims.UserID && !claims.Admin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		err = h.service.Delete(r.Context(), commID)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetProjectCommentsHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projIDStr := r.PathValue("proj_id")
		projID, err := strconv.Atoi(projIDStr)
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 10
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		offsetStr := r.URL.Query().Get("offset")
		offset := 0
		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 && o <= 10000 {
				offset = o
			} else if err == nil && o > 10000 {
				http.Error(w, "offset too large", http.StatusBadRequest)
				return
			}
		}

		comments, err := h.service.GetByProject(r.Context(), projID, limit, offset)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(comments)
	}
}
