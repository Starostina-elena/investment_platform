package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
)

var allowedDocMimeTypes = map[string]bool{
	"application/pdf":    true,
	"image/jpeg":         true,
	"image/png":          true,
	"image/tiff":         true,
	"application/msword": true, // .doc
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
	"application/vnd.ms-excel": true, // .xls
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true, // .xlsx
}

const MaxDocSize = 50 << 20 // 50 MB

func (s *service) UploadAvatar(ctx context.Context, orgID int, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	authorized, err := s.CheckUserOrgPermission(ctx, orgID, userID, "org_account_management")
	if err != nil {
		s.log.Error("failed to check user org permission", "error", err)
		return "", core.ErrNotAuthorized
	}
	if !authorized {
		return "", core.ErrNotAuthorized
	}

	_, err = s.Get(ctx, orgID)
	if err != nil {
		return "", core.ErrOrgNotFound
	}

	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w (supported: jpg, png, gif)", err)
	}
	s.log.Info("Decoded image", "format", format)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return "", fmt.Errorf("failed to encode to JPEG: %w", err)
	}

	objectName := fmt.Sprintf("orgpic_%d.jpg", orgID)

	imagePath, err := s.minio.PutAvatar(ctx, objectName, &buf, "image/jpeg")
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar to Minio: %w", err)
	}

	if err := s.repo.UpdateAvatarPath(ctx, orgID, &imagePath); err != nil {
		s.log.Error("failed to update avatar path", "error", err, "user_id", userID, "org_id", orgID)
		err2 := s.DeleteAvatar(ctx, orgID, userID, imagePath)
		if err2 != nil {
			s.log.Error("failed to delete avatar after update failure", "error", err2,
				"user_id", userID, "org_id", orgID, "path", imagePath)
		}
		s.log.Error("failed to update organisation avatar", "err", err)
		return "", err
	}

	return imagePath, nil
}

func (s *service) DeleteAvatar(ctx context.Context, orgID int, userID int, avatarPath string) error {
	org, err := s.Get(ctx, orgID)
	if err != nil {
		return core.ErrOrgNotFound
	}

	if org.AvatarPath == nil {
		return nil
	}

	if err := s.repo.UpdateAvatarPath(ctx, orgID, nil); err != nil {
		s.log.Error("failed to update avatar path", "error", err, "org_id", orgID, "user_id", userID)
		return err
	}

	err = s.minio.DeleteAvatar(ctx, *org.AvatarPath)

	if err != nil {
		s.log.Error("failed to delete avatar from storage", "error", err, "user_id", userID, "org_id", orgID, "path", *org.AvatarPath)
		if err2 := s.repo.UpdateAvatarPath(ctx, orgID, org.AvatarPath); err2 != nil {
			s.log.Error("failed to rollback avatar path after deletion failure", "error", err2, "user_id", userID, "org_id", orgID)
		}
		return err
	}

	return nil
}

func (s *service) UploadDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	authorized, err := s.CheckUserOrgPermission(ctx, orgID, userID, "org_account_management")
	if err != nil {
		s.log.Error("failed to check user org permission", "error", err)
		return "", core.ErrNotAuthorized
	}
	if !authorized {
		return "", core.ErrNotAuthorized
	}

	org, err := s.Get(ctx, orgID)
	if err != nil {
		return "", core.ErrOrgNotFound
	}

	if !docType.IsValidForOrgType(org.OrgType) {
		s.log.Warn("invalid doc type for org type", "user_id", userID, "org_id", orgID, "doc_type", docType, "org_type", org.OrgType)
		return "", fmt.Errorf("document type %s is not valid for organization type %s", docType, org.OrgType)
	}

	if fileHeader.Size > MaxDocSize {
		return "", fmt.Errorf("file is too large (max %d MB)", MaxDocSize>>20)
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if !allowedDocMimeTypes[contentType] {
		s.log.Warn("disallowed mime type for doc upload", "content_type", contentType, "user_id", userID, "org_id", orgID)
		return "", fmt.Errorf("file type not allowed: %s", contentType)
	}

	path, err := s.repo.GetDocPath(ctx, orgID, docType)
	if err != nil {
		return "", err
	}
	if path != "" {
		s.log.Info("overwriting existing doc", "org_id", orgID, "doc_type", docType)
		err = s.minio.DeleteDoc(ctx, path)
		if err != nil {
			s.log.Error("failed to delete existing doc", "error", err, "org_id", orgID, "doc_type", docType)
			return "", err
		}
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	objectName := fmt.Sprintf("org_%d_%s", orgID, docType)
	path, err = s.minio.PutDoc(ctx, objectName, &buf, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload doc: %w", err)
	}

	if err := s.repo.UpdateDocPath(ctx, orgID, docType, path); err != nil {
		s.log.Error("failed to update doc path", "error", err, "org_id", orgID, "doc_type", docType)
		_ = s.minio.DeleteDoc(ctx, path)
		return "", err
	}

	return path, nil
}

func (s *service) DeleteDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType) error {
	authorized, err := s.CheckUserOrgPermission(ctx, orgID, userID, "org_account_management")
	if err != nil {
		s.log.Error("failed to check user org permission", "error", err)
		return core.ErrNotAuthorized
	}
	if !authorized {
		return core.ErrNotAuthorized
	}

	org, err := s.Get(ctx, orgID)
	if err != nil {
		return core.ErrOrgNotFound
	}

	if !docType.IsValidForOrgType(org.OrgType) {
		s.log.Warn("invalid doc type for org type", "user_id", userID, "org_id", orgID, "doc_type", docType, "org_type", org.OrgType)
		return fmt.Errorf("document type %s is not valid for organization type %s", docType, org.OrgType)
	}

	path, err := s.repo.GetDocPath(ctx, orgID, docType)
	if err != nil {
		return err
	}
	if path != "" {
		err = s.repo.UpdateDocPath(ctx, orgID, docType, "")
		if err != nil {
			s.log.Error("failed to clear doc path in db", "error", err, "org_id", orgID, "doc_type", docType)
			return err
		}
		s.log.Info("deleting doc", "org_id", orgID, "doc_type", docType)
		err = s.minio.DeleteDoc(ctx, path)
		if err != nil {
			s.log.Error("failed to delete existing doc", "error", err, "org_id", orgID, "doc_type", docType)
			err2 := s.repo.UpdateDocPath(ctx, orgID, docType, path)
			if err2 != nil {
				s.log.Error("failed to restore doc path in db", "error", err2, "org_id", orgID, "doc_type", docType)
				return err2
			}
			return err
		}
	}

	return nil
}

func (s *service) DownloadDoc(ctx context.Context, orgID int, userID int, isAdmin bool, docType core.OrgDocType) ([]byte, string, error) {
	if !isAdmin {
		authorized, err := s.CheckUserOrgPermission(ctx, orgID, userID, "org_account_management")
		if err != nil {
			s.log.Error("failed to check user org permission", "error", err)
			return nil, "", core.ErrNotAuthorized
		}
		if !authorized {
			return nil, "", core.ErrNotAuthorized
		}
	}

	_, err := s.Get(ctx, orgID)
	if err != nil {
		return nil, "", core.ErrOrgNotFound
	}

	path, err := s.repo.GetDocPath(ctx, orgID, docType)
	if err != nil {
		return nil, "", err
	}
	if path == "" {
		return nil, "", core.ErrFileNotFound
	}

	data, contentType, err := s.minio.GetObject(ctx, path)
	if err != nil {
		return nil, "", err
	}

	return data, contentType, nil
}
