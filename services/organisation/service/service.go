package service

import (
	"context"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
	"github.com/Starostina-elena/investment_platform/services/organisation/repo"
	"github.com/Starostina-elena/investment_platform/services/organisation/storage"
)

type Service interface {
	Create(ctx context.Context, org core.Org) (*core.Org, error)
	Get(ctx context.Context, id int) (*core.Org, error)
	Update(ctx context.Context, org core.Org, userRequestedId int) (*core.Org, error)
	UploadAvatar(ctx context.Context, orgID int, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	DeleteAvatar(ctx context.Context, orgID int, userID int, avatarPath string) error
	UpdateAvatarPath(ctx context.Context, orgID int, avatarPath string) error
	UploadDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	DeleteDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType) error
	DownloadDoc(ctx context.Context, orgID int, userID int, isAdmin bool, docType core.OrgDocType) ([]byte, string, error)
}

type service struct {
	repo  repo.RepoInterface
	minio storage.MinioStorage
	log   slog.Logger
}

func NewService(r repo.RepoInterface, minio storage.MinioStorage, log slog.Logger) Service {
	return &service{repo: r, minio: minio, log: log}
}

func (s *service) Create(ctx context.Context, org core.Org) (*core.Org, error) {
	org.Balance = 0.0
	org.CreatedAt = time.Now()
	org.IsBanned = false
	id, err := s.repo.Create(ctx, &org)
	if err != nil {
		s.log.Error("failed to create organisation", "error", err)
		return nil, err
	}
	org.ID = id
	return &org, nil
}

func (s *service) Get(ctx context.Context, id int) (*core.Org, error) {
	org, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Error("failed to get organisation", "error", err)
		return nil, core.ErrOrgNotFound
	}

	org.SetIsRegistrationCompleted()

	return org, nil
}

func (s *service) Update(ctx context.Context, org core.Org, userRequestedId int) (*core.Org, error) {
	oldOrg, err := s.repo.Get(ctx, org.ID)
	if err != nil {
		s.log.Error("failed to get organisation for update", "error", err)
		return nil, core.ErrOrgNotFound
	}
	if oldOrg.OwnerId != userRequestedId {
		s.log.Error("user not authorized to update organisation", "org_id", org.ID, "user_id", userRequestedId)
		return nil, core.ErrNotAuthorized
	}

	org.OrgType = oldOrg.OrgType
	org.OrgTypeId = oldOrg.OrgTypeId
	org.Balance = oldOrg.Balance
	org.CreatedAt = oldOrg.CreatedAt
	org.IsBanned = oldOrg.IsBanned
	org.OwnerId = oldOrg.OwnerId
	org.SetIsRegistrationCompleted()

	updatedOrg, err := s.repo.Update(ctx, &org)
	if err != nil {
		s.log.Error("failed to update organisation", "error", err)
		return nil, err
	}

	return updatedOrg, nil
}

func (s *service) UpdateAvatarPath(ctx context.Context, orgID int, avatarPath string) error {
	var pathPtr *string
	if avatarPath != "" {
		pathPtr = &avatarPath
	}
	return s.repo.UpdateAvatarPath(ctx, orgID, pathPtr)
}
