package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
)

func (s *service) UploadAvatar(ctx context.Context, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
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

	objectName := fmt.Sprintf("userpic_%d.jpg", userID)

	imagePath, err := s.minio.PutObject(ctx, objectName, &buf, "image/jpeg")
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar to Minio: %w", err)
	}

	if err := s.UpdateAvatarPath(ctx, userID, imagePath); err != nil {
		s.log.Error("failed to update avatar path", "error", err, "user_id", userID)
		err2 := s.DeleteAvatar(ctx, userID, imagePath)
		if err2 != nil {
			s.log.Error("failed to delete avatar after update failure", "error", err2, "user_id", userID, "path", imagePath)
		}
		s.log.Error("failed to update user avatar", "err", err)
		return "", err
	}

	return imagePath, nil
}

func (s *service) DeleteAvatar(ctx context.Context, userID int, avatarPath string) error {
	if err := s.UpdateAvatarPath(ctx, userID, ""); err != nil {
		s.log.Error("failed to update avatar path", "error", err, "user_id", userID)
		return err
	}

	err := s.minio.DeleteAvatar(ctx, avatarPath)

	if err != nil {
		s.log.Error("failed to delete avatar from storage", "error", err, "user_id", userID, "path", avatarPath)
		if err2 := s.UpdateAvatarPath(ctx, userID, avatarPath); err2 != nil {
			s.log.Error("failed to rollback avatar path after deletion failure", "error", err2, "user_id", userID)
		}
		return err
	}

	return nil
}
