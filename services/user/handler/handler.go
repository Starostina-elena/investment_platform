package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/user/core"
	"github.com/Starostina-elena/investment_platform/services/user/service"
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler {
	return &Handler{service: s, log: log}
}

func CreateUserHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user core.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		u, err := h.service.Create(r.Context(), user)
		if err != nil {
			h.log.Error("failed to create user", "error", err)
			if err == core.ErrNicknameExists || err == core.ErrEmailExists {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, "Error while creating user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(u)
	}
}

func GetUserHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.log.Error("invalid user id", "id", idStr, "error", err)
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		u, err := h.service.Get(r.Context(), id)
		if err != nil {
			h.log.Error("failed to get user", "id", id, "error", err)
			http.Error(w, "No such user", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(u)
	}
}
