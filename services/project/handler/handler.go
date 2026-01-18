package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Starostina-elena/investment_platform/services/project/core"
	"github.com/Starostina-elena/investment_platform/services/project/middleware"
	"github.com/Starostina-elena/investment_platform/services/project/service"
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler {
	return &Handler{service: s, log: log}
}

type CreateProjectRequest struct {
	Name         string  `json:"name"`
	CreatorID    int     `json:"creator_id"`
	QuickPeek    string  `json:"quick_peek"`
	Content      string  `json:"content"`
	WantedMoney  float64 `json:"wanted_money"`
	DurationDays int     `json:"duration_days"`
}

func CreateProjectHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req CreateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode request", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.QuickPeek == "" || req.Content == "" || req.WantedMoney <= 0 || req.DurationDays <= 0 {
			http.Error(w, "Название, краткое и полное описание, срок и желаемая сумма обязательны", http.StatusBadRequest)
			return
		}
		if len(req.Name) > 128 || len(req.QuickPeek) > 128 || len(req.Content) > 1024 || req.WantedMoney > 1e9 || req.DurationDays > 36500 {
			http.Error(w, "Превышено максимальное количество символов в одном из полей", http.StatusBadRequest)
			return
		}

		p := core.Project{
			Name:         req.Name,
			CreatorID:    req.CreatorID,
			QuickPeek:    req.QuickPeek,
			Content:      req.Content,
			WantedMoney:  req.WantedMoney,
			DurationDays: req.DurationDays,
		}

		proj, err := h.service.Create(r.Context(), p, req.CreatorID, claims.UserID)
		if err != nil {
			h.log.Error("failed to create project", "error", err)
			switch err {
			case core.ErrInvalidInput:
				http.Error(w, "Некорректные входные данные", http.StatusBadRequest)
				return
			case core.ErrNotAuthorized:
				http.Error(w, "Нет прав для создания проекта в этой организации", http.StatusForbidden)
				return
			}
			http.Error(w, "Ошибка сервера при создании проекта", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(proj)
	}
}

func GetProjectHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.log.Error("invalid project id", "id", idStr, "error", err)
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		p, err := h.service.Get(r.Context(), id)
		if err != nil {
			h.log.Error("err while getting project", "id", id, "error", err)
			if err == core.ErrProjectNotFound {
				http.Error(w, "Проект не найден", http.StatusNotFound)
				return
			}
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(p)
	}
}

type UpdateProjectRequest struct {
	Name         string  `json:"name"`
	QuickPeek    string  `json:"quick_peek"`
	Content      string  `json:"content"`
	IsPublic     bool    `json:"is_public"`
	WantedMoney  float64 `json:"wanted_money"`
	DurationDays int     `json:"duration_days"`
}

func UpdateProjectHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.log.Error("invalid project id", "id", idStr, "error", err)
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		var req UpdateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode request", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.QuickPeek == "" || req.Content == "" || req.WantedMoney <= 0 || req.DurationDays <= 0 {
			http.Error(w, "Название, краткое и полное описание, срок и желаемая сумма обязательны", http.StatusBadRequest)
			return
		}
		if len(req.Name) > 128 || len(req.QuickPeek) > 128 || len(req.Content) > 1024 || req.WantedMoney > 1e9 || req.DurationDays > 36500 {
			http.Error(w, "Превышено максимальное количество символов в одном из полей", http.StatusBadRequest)
			return
		}

		p := core.Project{
			Name:         req.Name,
			QuickPeek:    req.QuickPeek,
			Content:      req.Content,
			IsPublic:     req.IsPublic,
			WantedMoney:  req.WantedMoney,
			DurationDays: req.DurationDays,
		}

		proj, err := h.service.Update(r.Context(), id, p, claims.UserID)
		if err != nil {
			h.log.Error("failed to update project", "error", err)
			switch err {
			case core.ErrProjectNotFound:
				http.Error(w, "Проект не найден", http.StatusNotFound)
				return
			case core.ErrNotAuthorized:
				http.Error(w, "Нет прав для изменения проекта", http.StatusForbidden)
				return
			case core.ErrInvalidInput:
				http.Error(w, "Некорректные входные данные", http.StatusBadRequest)
				return
			}
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(proj)
	}
}

func GetProjectListHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit := 10
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}
		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil {
				offset = o
			}
		}

		projects, err := h.service.GetList(r.Context(), limit, offset)
		if err != nil {
			h.log.Error("failed to get projects list", "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(projects)
	}
}

func GetProjectsByCreatorHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creatorIDStr := r.PathValue("creator_id")
		creatorID, err := strconv.Atoi(creatorIDStr)
		if err != nil {
			h.log.Error("invalid creator id", "id", creatorIDStr, "error", err)
			http.Error(w, "Некорректный id создателя", http.StatusBadRequest)
			return
		}

		projects, err := h.service.GetByCreator(r.Context(), creatorID)
		if err != nil {
			h.log.Error("failed to get projects by creator", "creator_id", creatorID, "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(projects)
	}
}

func BanProjectHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !claims.Admin || claims.Banned {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		projectIDStr := r.PathValue("id")
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			h.log.Error("invalid project id", "id", projectIDStr, "error", err)
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		banStr := strings.TrimSpace(r.URL.Query().Get("ban"))
		ban, err := strconv.ParseBool(banStr)
		if err != nil {
			h.log.Error("invalid ban value", "value", banStr, "error", err)
			http.Error(w, "Некорректное значение ban", http.StatusBadRequest)
			return
		}

		err = h.service.BanProject(r.Context(), projectID, ban)
		if err != nil {
			h.log.Error("failed to ban/unban project", "project_id", projectID, "banned", ban, "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func MarkProjectCompletedHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		projectIDStr := r.PathValue("id")
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			h.log.Error("invalid project id", "id", projectIDStr, "error", err)
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		completedStr := strings.TrimSpace(r.URL.Query().Get("completed"))
		completed, err := strconv.ParseBool(completedStr)
		if err != nil {
			h.log.Error("invalid completed value", "value", completedStr, "error", err)
			http.Error(w, "Некорректное значение completed", http.StatusBadRequest)
			return
		}

		err = h.service.MarkProjectCompleted(r.Context(), projectID, claims.UserID, completed)
		if err != nil {
			h.log.Error("failed to mark project as completed", "project_id", projectID, "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
