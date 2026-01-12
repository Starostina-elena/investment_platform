package handler

import (
	"net/http"
	"strconv"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
	"github.com/Starostina-elena/investment_platform/services/organisation/middleware"
)

type DocUploadResponse struct {
	Path string `json:"path"`
}

func UploadDocHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		orgIdStr := r.PathValue("org_id")
		orgID, err := strconv.Atoi(orgIdStr)
		if err != nil {
			h.log.Error("invalid org id", "id", orgIdStr, "error", err)
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		docType := core.OrgDocType(r.PathValue("doc_type"))
		if !core.IsValidDocType(docType) {
			h.log.Warn("invalid doc type attempted", "doc_type", docType, "user_id", claims.UserID, "org_id", orgID)
			http.Error(w, "Некорректный тип документа", http.StatusBadRequest)
			return
		}

		err = r.ParseMultipartForm(51 << 20)
		if err != nil {
			h.log.Error("failed to parse multipart form", "error", err)
			http.Error(w, "Файл слишком большой", http.StatusRequestEntityTooLarge)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			h.log.Error("failed to get file from form", "error", err)
			http.Error(w, "Файл не предоставлен", http.StatusBadRequest)
			return
		}
		defer file.Close()

		path, err := h.service.UploadDoc(r.Context(), orgID, claims.UserID, docType, file, fileHeader)
		if err != nil {
			h.log.Error("failed to upload org doc", "error", err, "org_id", orgID, "doc_type", docType)
			switch err {
			case core.ErrOrgNotFound:
				http.Error(w, "Организация не найдена", http.StatusNotFound)
				return
			case core.ErrNotAuthorized:
				http.Error(w, "Недостаточно прав", http.StatusForbidden)
				return
			}
			http.Error(w, "Ошибка при загрузке документа", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"path":"` + path + `"}`))
	}
}

func DeleteDocHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		orgIdStr := r.PathValue("org_id")
		orgID, err := strconv.Atoi(orgIdStr)
		if err != nil {
			h.log.Error("invalid org id", "id", orgIdStr, "error", err)
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		docType := core.OrgDocType(r.PathValue("doc_type"))
		if !core.IsValidDocType(docType) {
			h.log.Warn("invalid doc type attempted", "doc_type", docType, "user_id", claims.UserID, "org_id", orgID)
			http.Error(w, "Некорректный тип документа", http.StatusBadRequest)
			return
		}

		err = h.service.DeleteDoc(r.Context(), orgID, claims.UserID, docType)
		if err != nil {
			h.log.Error("failed to delete org doc", "error", err, "org_id", orgID, "doc_type", docType)
			switch err {
			case core.ErrOrgNotFound:
				http.Error(w, "Организация не найдена", http.StatusNotFound)
				return
			case core.ErrNotAuthorized:
				http.Error(w, "Нет доступа", http.StatusForbidden)
				return
			}
			http.Error(w, "Ошибка при удалении документа", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DownloadDocHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.FromContext(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		orgIdStr := r.PathValue("org_id")
		orgID, err := strconv.Atoi(orgIdStr)
		if err != nil {
			h.log.Error("invalid org id", "id", orgIdStr, "error", err)
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		docType := core.OrgDocType(r.PathValue("doc_type"))
		if !core.IsValidDocType(docType) {
			h.log.Warn("invalid doc type attempted", "doc_type", docType, "user_id", claims.UserID, "org_id", orgID)
			http.Error(w, "Некорректный тип документа", http.StatusBadRequest)
			return
		}

		data, contentType, err := h.service.DownloadDoc(r.Context(), orgID, claims.UserID, claims.Admin, docType)
		if err != nil {
			h.log.Error("failed to download org doc", "error", err, "org_id", orgID, "doc_type", docType)
			switch err {
			case core.ErrOrgNotFound:
				http.Error(w, "Организация не найдена", http.StatusNotFound)
				return
			case core.ErrNotAuthorized:
				http.Error(w, "Нет доступа", http.StatusForbidden)
				return
			case core.ErrFileNotFound:
				http.Error(w, "Файл не найден", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to download doc", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}
