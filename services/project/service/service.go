package service

import (
	"context"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/Starostina-elena/investment_platform/services/project/clients"
	"github.com/Starostina-elena/investment_platform/services/project/core"
	"github.com/Starostina-elena/investment_platform/services/project/repo"
	"github.com/Starostina-elena/investment_platform/services/project/storage"
)

type Service interface {
	Create(ctx context.Context, req core.Project, creatorID int, userID int) (*core.Project, error)
	Get(ctx context.Context, id int) (*core.Project, error)
	Update(ctx context.Context, projectID int, p core.Project, userID int) (*core.Project, error)
	GetList(ctx context.Context, limit, offset int) ([]core.Project, error)
	GetByCreator(ctx context.Context, creatorID int) ([]core.Project, error)
	UpdatePicturePath(ctx context.Context, projectID int, picturePath string) error
	BanProject(ctx context.Context, projectID int, banned bool) error
	MarkProjectCompleted(ctx context.Context, projectID int, userID int, completed bool) error
	UploadPicture(ctx context.Context, projectID int, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	deletePicture(ctx context.Context, projectID int, picturePath string) error
	DeletePictureFromProject(ctx context.Context, projectID int, userID int) error
}

type service struct {
	repo      repo.RepoInterface
	orgClient *clients.OrgClient
	minio     *storage.MinioStorage
	log       slog.Logger
}

func NewService(r repo.RepoInterface, orgClient *clients.OrgClient, minioStorage *storage.MinioStorage, log slog.Logger) Service {
	return &service{repo: r, orgClient: orgClient, minio: minioStorage, log: log}
}

func (s *service) Create(ctx context.Context, p core.Project, creatorID int, userID int) (*core.Project, error) {
	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, creatorID, userID, "project_management")
	if err != nil {
		s.log.Error("failed to check organisation permission", "error", err)
		return nil, err
	}
	if !allowed {
		return nil, core.ErrNotAuthorized
	}

	p.IsPublic = true
	p.IsCompleted = false
	p.CurrentMoney = 0.0
	p.CreatedAt = time.Now()
	p.IsBanned = false

	id, err := s.repo.Create(ctx, &p)
	if err != nil {
		s.log.Error("failed to create project", "error", err)
		return nil, err
	}
	p.ID = id
	return &p, nil
}

func (s *service) Get(ctx context.Context, id int) (*core.Project, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, projectID int, p core.Project, userID int) (*core.Project, error) {
	existingProject, err := s.repo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}

	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, existingProject.CreatorID, userID, "project_management")
	if err != nil || !allowed {
		return nil, core.ErrNotAuthorized
	}

	existingProject.Name = p.Name
	existingProject.QuickPeek = p.QuickPeek
	existingProject.Content = p.Content
	existingProject.IsPublic = p.IsPublic
	existingProject.WantedMoney = p.WantedMoney
	existingProject.DurationDays = p.DurationDays

	updatedProject, err := s.repo.Update(ctx, existingProject)
	if err != nil {
		s.log.Error("failed to update project", "error", err)
		return nil, err
	}
	return updatedProject, nil
}

func (s *service) GetList(ctx context.Context, limit, offset int) ([]core.Project, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetList(ctx, limit, offset)
}

func (s *service) GetByCreator(ctx context.Context, creatorID int) ([]core.Project, error) {
	return s.repo.GetByCreator(ctx, creatorID)
}

func (s *service) UpdatePicturePath(ctx context.Context, projectID int, picturePath string) error {
	return s.repo.UpdatePicturePath(ctx, projectID, &picturePath)
}

func (s *service) BanProject(ctx context.Context, projectID int, banned bool) error {
	return s.repo.BanProject(ctx, projectID, banned)
}

func (s *service) MarkProjectCompleted(ctx context.Context, projectID int, userID int, completed bool) error {
	existingProject, err := s.repo.Get(ctx, projectID)
	if err != nil {
		return err
	}
	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, existingProject.CreatorID, userID, "project_management")
	if err != nil || !allowed {
		return core.ErrNotAuthorized
	}
	return s.repo.MarkProjectCompleted(ctx, projectID, completed)
}
