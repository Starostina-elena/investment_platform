package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
	"github.com/Starostina-elena/investment_platform/services/organisation/middleware"
)

func UploadAvatarHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		userID := claims.UserID

		orgIdStr := r.PathValue("org_id")
		orgId, err := strconv.Atoi(orgIdStr)
		if err != nil {
			h.log.Error("invalid org id", "id", orgIdStr, "error", err)
			http.Error(w, "Некорретный id", http.StatusBadRequest)
			return
		}

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			h.log.Error("failed to parse multipart form", "error", err)
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("avatar")
		if err != nil {
			h.log.Error("failed to get file from form", "error", err)
			http.Error(w, "no file provided", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if fileHeader.Size > 5<<20 {
			http.Error(w, "file too large (max 5MB)", http.StatusBadRequest)
			return
		}

		avatarName, err := h.service.UploadAvatar(r.Context(), orgId, userID, file, fileHeader)
		if err != nil {
			h.log.Error("failed to upload avatar", "error", err, "user_id", userID, "org_id", orgId)
			if err == core.ErrOrgNotFound {
				http.Error(w, "organisation not found", http.StatusNotFound)
				return
			}
			if err == core.ErrNotAuthorized {
				http.Error(w, "Нет прав для загрузки аватара", http.StatusForbidden)
				return
			}
			http.Error(w, "failed to upload avatar", http.StatusInternalServerError)
			return
		}

		h.log.Info("avatar uploaded successfully", "user_id", userID, "org_id", orgId, "filename", avatarName)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":     "avatar uploaded successfully",
			"avatar_path": avatarName,
		})
	}
}

func DeleteAvatarHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		userID := claims.UserID

		orgIdStr := r.PathValue("org_id")
		orgId, err := strconv.Atoi(orgIdStr)
		if err != nil {
			h.log.Error("invalid org id", "id", orgIdStr, "error", err)
			http.Error(w, "Некорретный id", http.StatusBadRequest)
			return
		}

		org, err := h.service.Get(r.Context(), orgId)
		if err != nil {
			h.log.Error("org not found", "id", orgId, "error", err)
			http.Error(w, "Организация не найдена", http.StatusNotFound)
			return
		}
		if org.OwnerId != userID {
			h.log.Error("user not authorized to delete avatar", "org_id", orgId, "user_id", userID)
			http.Error(w, "Нет прав для удаления аватара", http.StatusForbidden)
			return
		}
		if org.AvatarPath == nil {
			http.Error(w, "no avatar to delete", http.StatusNotFound)
			return
		}

		if err := h.service.DeleteAvatar(r.Context(), orgId, userID, *org.AvatarPath); err != nil {
			h.log.Error("failed to delete avatar from storage", "error", err, "path", *org.AvatarPath)
			http.Error(w, "failed to delete avatar", http.StatusInternalServerError)
			return
		}

		h.log.Info("avatar deleted successfully", "user_id", userID, "org_id", orgId)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "avatar deleted successfully",
		})
	}
}
