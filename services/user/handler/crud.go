package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/user/core"
	"github.com/Starostina-elena/investment_platform/services/user/middleware"
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
			http.Error(w, "Ошибка сервера при создании пользователя", http.StatusInternalServerError)
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
			http.Error(w, "Некорретный id", http.StatusBadRequest)
			return
		}

		userRequested := middleware.FromContext(r.Context())
		if !userRequested.Admin && userRequested.UserID != id {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		u, err := h.service.Get(r.Context(), id)
		if err != nil {
			h.log.Error("failed to get user", "id", id, "error", err)
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(u)
	}
}

func UpdateUserHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uc := middleware.FromContext(r.Context())
		if uc == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var user core.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		user.ID = uc.UserID
		u, err := h.service.Update(r.Context(), user)
		if err != nil {
			h.log.Error("failed to update user", "id", user.ID, "error", err)
			h.log.Error("failed to create user", "error", err)
			if err == core.ErrNicknameExists || err == core.ErrEmailExists {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, "Ошибка сервера при обновлении пользователя", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(u)
	}
}
