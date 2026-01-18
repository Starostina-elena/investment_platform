package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"

	"github.com/Starostina-elena/investment_platform/services/project/core"
)

func (s *service) UploadPicture(ctx context.Context, projectID int, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w (supported: jpg, png, gif)", err)
	}
	s.log.Info("Decoded image", "format", format)

	existingProject, err := s.repo.Get(ctx, projectID)
	if err != nil {
		return "", err
	}
	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, existingProject.CreatorID, userID, "project_management")
	if err != nil || !allowed {
		return "", core.ErrNotAuthorized
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return "", fmt.Errorf("failed to encode to JPEG: %w", err)
	}

	objectName := fmt.Sprintf("projpic_%d.jpg", projectID)

	imagePath, err := s.minio.PutObject(ctx, objectName, &buf, "image/jpeg")
	if err != nil {
		return "", fmt.Errorf("failed to upload picture to Minio: %w", err)
	}

	if err := s.repo.UpdatePicturePath(ctx, projectID, &imagePath); err != nil {
		s.log.Error("failed to update picture path", "error", err, "project_id", projectID)
		err2 := s.deletePicture(ctx, projectID, imagePath)
		if err2 != nil {
			s.log.Error("failed to delete picture after update failure", "error", err2, "project_id", projectID, "path", imagePath)
		}
		s.log.Error("failed to update project picture", "err", err)
		return "", err
	}

	return imagePath, nil
}

func (s *service) DeletePictureFromProject(ctx context.Context, projectID int, userID int) error {
	existingProject, err := s.repo.Get(ctx, projectID)
	if err != nil {
		return err
	}
	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, existingProject.CreatorID, userID, "project_management")
	if err != nil {
		s.log.Error("failed to check permission", "error", err, "user_id", userID, "org_id", existingProject.CreatorID)
		return err
	}
	if !allowed {
		return core.ErrNotAuthorized
	}

	if existingProject.QuickPeekPicturePath == nil {
		return nil
	}

	return s.deletePicture(ctx, projectID, *existingProject.QuickPeekPicturePath)
}

func (s *service) deletePicture(ctx context.Context, projectID int, picturePath string) error {
	if err := s.repo.UpdatePicturePath(ctx, projectID, nil); err != nil {
		s.log.Error("failed to update picture path", "error", err, "project_id", projectID)
		return err
	}

	err := s.minio.DeletePicture(ctx, picturePath)

	if err != nil {
		s.log.Error("failed to delete picture from storage", "error", err, "project_id", projectID, "path", picturePath)
		if err2 := s.repo.UpdatePicturePath(ctx, projectID, &picturePath); err2 != nil {
			s.log.Error("failed to rollback picture path after deletion failure", "error", err2, "project_id", projectID)
		}
		return err
	}

	return nil
}
