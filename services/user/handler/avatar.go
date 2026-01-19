package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Starostina-elena/investment_platform/services/user/middleware"
)

func UploadAvatarHandler(h *Handler) http.HandlerFunc {
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
		userID := claims.UserID

		err := r.ParseMultipartForm(10 << 20)
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

		avatarName, err := h.service.UploadAvatar(r.Context(), userID, file, fileHeader)
		if err != nil {
			h.log.Error("failed to upload avatar", "error", err, "user_id", userID)
			http.Error(w, "failed to upload avatar", http.StatusInternalServerError)
			return
		}

		h.log.Info("avatar uploaded successfully", "user_id", userID, "filename", avatarName)

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

		user, err := h.service.Get(r.Context(), userID)
		if err != nil {
			h.log.Error("failed to get user", "error", err, "user_id", userID)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		if user.AvatarPath == nil || *user.AvatarPath == "" {
			http.Error(w, "user has no avatar", http.StatusNotFound)
			return
		}

		if err := h.service.DeleteAvatar(r.Context(), userID, *user.AvatarPath); err != nil {
			h.log.Error("failed to delete avatar from storage", "error", err, "path", *user.AvatarPath)
			http.Error(w, "failed to delete avatar", http.StatusInternalServerError)
			return
		}

		h.log.Info("avatar deleted successfully", "user_id", userID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "avatar deleted successfully",
		})
	}
}
