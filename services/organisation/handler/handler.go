package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
	"github.com/Starostina-elena/investment_platform/services/organisation/middleware"
	"github.com/Starostina-elena/investment_platform/services/organisation/service"
)

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(s service.Service, log slog.Logger) *Handler {
	return &Handler{service: s, log: log}
}

type CreateOrgRequest struct {
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	OrgType  core.OrgType   `json:"org_type"`
	PhysFace *core.PhysFace `json:"phys_face,omitempty"`
	JurFace  *core.JurFace  `json:"jur_face,omitempty"`
	IPFace   *core.IPFace   `json:"ip_face,omitempty"`
}

func CreateOrgHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req CreateOrgRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode request", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Email == "" || req.OrgType == "" {
			http.Error(w, "Название, электронная почта и тип организации обязательны", http.StatusBadRequest)
			return
		}

		switch req.OrgType {
		case core.OrgTypeIP:
			if req.IPFace == nil {
				http.Error(w, "Поля ip_face обязательны для типа организации ip", http.StatusBadRequest)
				return
			}
		case core.OrgTypePhys:
			if req.PhysFace == nil {
				http.Error(w, "Поля phys_face обязательны для типа организации phys", http.StatusBadRequest)
				return
			}
		case core.OrgTypeJur:
			if req.JurFace == nil {
				http.Error(w, "Поля jur_face обязательны для типа организации jur", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "invalid org_type", http.StatusBadRequest)
			return
		}

		org := core.Org{
			OrgBase: core.OrgBase{
				Name:    req.Name,
				Email:   req.Email,
				OrgType: req.OrgType,
				OwnerId: claims.UserID,
			},
			PhysFace: req.PhysFace,
			JurFace:  req.JurFace,
			IPFace:   req.IPFace,
		}

		o, err := h.service.Create(r.Context(), org)
		if err != nil {
			h.log.Error("failed to create org", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
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
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}
		o, err := h.service.Get(r.Context(), id)
		if err != nil {
			h.log.Error("org not found", "id", id, "error", err)
			http.Error(w, "Организация не найдена", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(o)
	}
}
