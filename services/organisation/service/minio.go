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

func (s *service) UploadAvatar(ctx context.Context, orgID int, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	org, err := s.Get(ctx, orgID)
	if err != nil {
		return "", core.ErrOrgNotFound
	}
	if org.OwnerId != userID {
		s.log.Error("user is not the creator of the organisation", "user_id", userID, "org_id", orgID)
		return "", core.ErrNotAuthorized
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
	if org.OwnerId != userID {
		s.log.Error("user is not the creator of the organisation", "user_id", userID, "org_id", orgID)
		return core.ErrNotAuthorized
	}

	if err := s.repo.UpdateAvatarPath(ctx, orgID, nil); err != nil {
		s.log.Error("failed to update avatar path", "error", err, "org_id", orgID, "user_id", userID)
		return err
	}

	err = s.minio.DeleteAvatar(ctx, avatarPath)

	if err != nil {
		s.log.Error("failed to delete avatar from storage", "error", err, "user_id", userID, "org_id", orgID, "path", avatarPath)
		if err2 := s.repo.UpdateAvatarPath(ctx, orgID, &avatarPath); err2 != nil {
			s.log.Error("failed to rollback avatar path after deletion failure", "error", err2, "user_id", userID, "org_id", orgID)
		}
		return err
	}

	return nil
}
