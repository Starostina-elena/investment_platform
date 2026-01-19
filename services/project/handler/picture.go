package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/project/core"
	"github.com/Starostina-elena/investment_platform/services/project/middleware"
)

func UploadPictureHandler(h *Handler) http.HandlerFunc {
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
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			h.log.Error("failed to parse multipart form", "error", err)
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("picture")
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

		allowedMimes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
		}
		if !allowedMimes[fileHeader.Header.Get("Content-Type")] {
			http.Error(w, "invalid file type, allowed: jpg, png, gif", http.StatusBadRequest)
			return
		}

		picturePath, err := h.service.UploadPicture(r.Context(), projectID, claims.UserID, file, fileHeader)
		if err != nil {
			if err == core.ErrNotAuthorized {
				h.log.Warn("unauthorized upload attempt", "project_id", projectID, "user_id", claims.UserID)
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}
			h.log.Error("failed to upload picture", "error", err, "project_id", projectID)
			http.Error(w, "failed to upload picture", http.StatusInternalServerError)
			return
		}

		h.log.Info("picture uploaded successfully", "project_id", projectID, "filename", picturePath)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "picture uploaded successfully",
			"picture_path": picturePath,
		})
	}
}

func DeletePictureHandler(h *Handler) http.HandlerFunc {
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
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		err = h.service.DeletePictureFromProject(r.Context(), projectID, claims.UserID)
		if err != nil {
			if err == core.ErrNotAuthorized {
				h.log.Warn("unauthorized delete attempt", "project_id", projectID, "user_id", claims.UserID)
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}
			h.log.Error("failed to delete picture", "error", err, "project_id", projectID)
			http.Error(w, "failed to delete picture", http.StatusInternalServerError)
			return
		}

		h.log.Info("picture deleted successfully", "project_id", projectID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "picture deleted successfully",
		})
	}
}
