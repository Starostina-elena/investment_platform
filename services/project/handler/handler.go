package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Starostina-elena/investment_platform/services/project/core"
	"github.com/Starostina-elena/investment_platform/services/project/middleware"
	"github.com/Starostina-elena/investment_platform/services/project/service"
)

const (
	MaxNameLength      = 128
	MaxQuickPeekLength = 128
	MaxContentLength   = 1024
	MaxWantedMoney     = 1e9
	MaxDurationDays    = 36500
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler {
	return &Handler{service: s, log: log}
}

func ValidateProjectInput(req interface{}) error {
	switch r := req.(type) {
	case CreateProjectRequest:
		if r.Name == "" || r.QuickPeek == "" || r.Content == "" || r.WantedMoney <= 0 || r.DurationDays <= 0 {
			return fmt.Errorf("название, краткое и полное описание, срок и желаемая сумма обязательны")
		}
		if len(r.Name) > MaxNameLength {
			return fmt.Errorf("название слишком длинное (макс %d символов)", MaxNameLength)
		}
		if len(r.QuickPeek) > MaxQuickPeekLength {
			return fmt.Errorf("краткое описание слишком длинное (макс %d символов)", MaxQuickPeekLength)
		}
		if len(r.Content) > MaxContentLength {
			return fmt.Errorf("полное описание слишком длинное (макс %d символов)", MaxContentLength)
		}
		if r.WantedMoney > MaxWantedMoney {
			return fmt.Errorf("желаемая сумма слишком велика (макс %.0f)", MaxWantedMoney)
		}
		if r.DurationDays > MaxDurationDays {
			return fmt.Errorf("срок слишком велик (макс %d дней)", MaxDurationDays)
		}
	case UpdateProjectRequest:
		if r.Name == "" || r.QuickPeek == "" || r.Content == "" || r.WantedMoney <= 0 || r.DurationDays <= 0 {
			return fmt.Errorf("название, краткое и полное описание, срок и желаемая сумма обязательны")
		}
		if len(r.Name) > MaxNameLength {
			return fmt.Errorf("название слишком длинное (макс %d символов)", MaxNameLength)
		}
		if len(r.QuickPeek) > MaxQuickPeekLength {
			return fmt.Errorf("краткое описание слишком длинное (макс %d символов)", MaxQuickPeekLength)
		}
		if len(r.Content) > MaxContentLength {
			return fmt.Errorf("полное описание слишком длинное (макс %d символов)", MaxContentLength)
		}
		if r.WantedMoney > MaxWantedMoney {
			return fmt.Errorf("желаемая сумма слишком велика (макс %.0f)", MaxWantedMoney)
		}
		if r.DurationDays > MaxDurationDays {
			return fmt.Errorf("срок слишком велик (макс %d дней)", MaxDurationDays)
		}
	}
	return nil
}

type CreateProjectRequest struct {
	Name             string  `json:"name"`
	CreatorID        int     `json:"creator_id"`
	QuickPeek        string  `json:"quick_peek"`
	Content          string  `json:"content"`
	WantedMoney      float64 `json:"wanted_money"`
	DurationDays     int     `json:"duration_days"`
	MonetizationType string  `json:"monetization_type"`
	Percent          float64 `json:"percent"`
}

func CreateProjectHandler(h *Handler) http.HandlerFunc {
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

		var req CreateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode request", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if err := ValidateProjectInput(req); err != nil {
			h.log.Warn("invalid project input", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.MonetizationType == "fixed_percent" || req.MonetizationType == "time_percent" {
			if req.Percent <= 0 {
				req.Percent = 5.0
			}
		} else {
			req.Percent = 0.0
		}

		p := core.Project{
			Name:             req.Name,
			CreatorID:        req.CreatorID,
			QuickPeek:        req.QuickPeek,
			Content:          req.Content,
			WantedMoney:      req.WantedMoney,
			DurationDays:     req.DurationDays,
			MonetizationType: req.MonetizationType,
			Percent:          req.Percent,
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
		if claims.Banned {
			http.Error(w, "forbidden", http.StatusForbidden)
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

		if err := ValidateProjectInput(req); err != nil {
			h.log.Warn("invalid project input", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		monetizationType := r.URL.Query().Get("type")

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

		if monetizationType != "" {
			validTypes := map[string]bool{
				"charity":       true,
				"custom":        true,
				"fixed_percent": true,
				"time_percent":  true,
			}
			if !validTypes[monetizationType] {
				http.Error(w, "Некорректный тип монетизации", http.StatusBadRequest)
				return
			}
		}

		projects, err := h.service.GetList(r.Context(), limit, offset, monetizationType)
		if err != nil {
			h.log.Error("failed to get projects list", "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(projects)
	}
}

func GetPublicProjectsByCreatorHandler(h *Handler) http.HandlerFunc {
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

func GetAllProjectsByCreatorHandler(h *Handler) http.HandlerFunc {
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
		creatorIDStr := r.PathValue("creator_id")
		creatorID, err := strconv.Atoi(creatorIDStr)
		if err != nil {
			h.log.Error("invalid creator id", "id", creatorIDStr, "error", err)
			http.Error(w, "Некорректный id создателя", http.StatusBadRequest)
			return
		}

		projects, err := h.service.GetAllByCreator(r.Context(), creatorID, claims.UserID, claims.Admin)
		if err != nil {
			if err == core.ErrNotAuthorized {
				h.log.Warn("unauthorized access attempt", "creator_id", creatorID, "user_id", claims.UserID)
				http.Error(w, "Нет прав для просмотра всех проектов этой организации", http.StatusForbidden)
				return
			}
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

func ChangeProjectPublicityHandler(h *Handler) http.HandlerFunc {
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

		projectIDStr := r.PathValue("id")
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			h.log.Error("invalid project id", "id", projectIDStr, "error", err)
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		publicStr := strings.TrimSpace(r.URL.Query().Get("public"))
		isPublic, err := strconv.ParseBool(publicStr)
		if err != nil {
			h.log.Error("invalid public value", "value", publicStr, "error", err)
			http.Error(w, "Некорректное значение public", http.StatusBadRequest)
			return
		}

		err = h.service.ChangeProjectPublicity(r.Context(), projectID, claims.UserID, isPublic)
		if err != nil {
			if err == core.ErrNotAuthorized {
				h.log.Warn("unauthorized attempt to change project publicity", "project_id", projectID, "user_id", claims.UserID)
				http.Error(w, "Нет прав для изменения публичности проекта", http.StatusForbidden)
				return
			}
			h.log.Error("failed to change project publicity", "project_id", projectID, "error", err)
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
		if claims.Banned {
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

		completedStr := strings.TrimSpace(r.URL.Query().Get("completed"))
		completed, err := strconv.ParseBool(completedStr)
		if err != nil {
			h.log.Error("invalid completed value", "value", completedStr, "error", err)
			http.Error(w, "Некорректное значение completed", http.StatusBadRequest)
			return
		}

		err = h.service.MarkProjectCompleted(r.Context(), projectID, claims.UserID, completed)
		if err != nil {
			if err == core.ErrNotAuthorized {
				h.log.Warn("unauthorized attempt to mark project completed", "project_id", projectID, "user_id", claims.UserID)
				http.Error(w, "Нет прав для изменения статуса проекта", http.StatusForbidden)
				return
			}
			if err == core.ErrPaybackStarted {
				h.log.Warn("attempt to change is_completed when payback has started", "project_id", projectID)
				http.Error(w, "Невозможно изменить статус завершенности проекта после начала возврата средств инвесторам", http.StatusBadRequest)
				return
			}
			h.log.Error("failed to mark project as completed", "project_id", projectID, "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func StartPaybackHandler(h *Handler) http.HandlerFunc {
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

		err = h.service.StartPayback(r.Context(), projectID, claims.UserID)
		if err != nil {
			if err == core.ErrNotAuthorized {
				h.log.Warn("unauthorized attempt to start payback", "project_id", projectID, "user_id", claims.UserID)
				http.Error(w, "Нет прав для запуска возврата средств (требуется money_management)", http.StatusForbidden)
				return
			}
			if err == core.ErrPaybackStarted {
				h.log.Warn("payback already started", "project_id", projectID)
				http.Error(w, "Возврат средств уже запущен для этого проекта", http.StatusBadRequest)
				return
			}
			if err == core.ErrPaybackNotSupported {
				h.log.Warn("payback not supported for monetization type", "project_id", projectID)
				http.Error(w, "Возврат средств не поддерживается для проектов с типом монетизации charity или custom", http.StatusBadRequest)
				return
			}
			if err == core.ErrNotEnoughFunds {
				h.log.Warn("not enough funds to complete payback", "project_id", projectID)
				http.Error(w, "Недостаточно средств на проекте для полного возврата инвесторам", http.StatusBadRequest)
				return
			}
			h.log.Error("failed to start payback", "project_id", projectID, "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type AddFundsRequest struct {
	Amount float64 `json:"amount"`
}

func AddFundsHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var req AddFundsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Amount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}

		err = h.service.AddFunds(r.Context(), id, req.Amount)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
