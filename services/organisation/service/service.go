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
	GetPublicInfoOrg(ctx context.Context, id int) (*core.Org, error)
	Update(ctx context.Context, org core.Org) (*core.Org, error)
	UploadAvatar(ctx context.Context, orgID int, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	DeleteAvatar(ctx context.Context, orgID int, userID int, avatarPath string) error
	UpdateAvatarPath(ctx context.Context, orgID int, avatarPath string) error
	UploadDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	DeleteDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType) error
	DownloadDoc(ctx context.Context, orgID int, userID int, isAdmin bool, docType core.OrgDocType) ([]byte, string, error)
	GetUsersOrgs(ctx context.Context, userID int) ([]core.Org, error)
	BanOrg(ctx context.Context, orgID int, banned bool) error
	CheckUserOrgPermission(ctx context.Context, orgID int, userID int, permission core.OrgPermission) (bool, error)
	AddEmployee(ctx context.Context, orgID int, userRequested int, userID int, orgAccMgmt, moneyMgmt, projMgmt bool) error
	GetOrgEmployees(ctx context.Context, orgID int) ([]core.OrgEmployee, error)
	UpdateEmployeePermissions(ctx context.Context, orgID int, userRequested int, userID int, orgAccMgmt, moneyMgmt, projMgmt bool) error
	DeleteEmployee(ctx context.Context, orgID int, userRequested int, userID int) error
	TransferOwnership(ctx context.Context, orgID int, userRequested int, newOwnerID int) error
	ChangeBalance(ctx context.Context, orgID int, delta float64) error
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

func (s *service) GetPublicInfoOrg(ctx context.Context, id int) (*core.Org, error) {
	org, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Error("failed to get public info of organisation", "error", err)
		return nil, core.ErrOrgNotFound
	}

	org.SetIsRegistrationCompleted()

	org.Balance = 0.0
	switch org.OrgType {
	case core.OrgTypePhys:
		org.PhysFace.INN = ""
		org.PhysFace.PassportSeries = 0
		org.PhysFace.PassportNumber = 0
		org.PhysFace.PassportGivenBy = ""
		org.PhysFace.PostAddress = ""
		org.PhysFace.RegistrationAddress = ""
	case core.OrgTypeJur:
		org.JurFace.INN = ""
		org.JurFace.KPP = ""
		org.JurFace.JurAddress = ""
		org.JurFace.FactAddress = ""
		org.JurFace.PostAddress = ""
	case core.OrgTypeIP:
		org.IPFace.IpSvidSerial = 0
		org.IPFace.IpSvidNumber = 0
		org.IPFace.IpSvidGivenBy = ""
		org.IPFace.INN = ""
		org.IPFace.JurAddress = ""
		org.IPFace.FactAddress = ""
		org.IPFace.PostAddress = ""
	}

	return org, nil
}

func (s *service) Update(ctx context.Context, org core.Org) (*core.Org, error) {
	oldOrg, err := s.repo.Get(ctx, org.ID)
	if err != nil {
		s.log.Error("failed to get organisation for update", "error", err)
		return nil, core.ErrOrgNotFound
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

func (s *service) GetUsersOrgs(ctx context.Context, userID int) ([]core.Org, error) {
	orgs, err := s.repo.GetUsersOrgs(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user's organisations", "error", err, "user_id", userID)
		return nil, err
	}

	for i := range orgs {
		orgs[i].SetIsRegistrationCompleted()
	}

	return orgs, nil
}

func (s *service) BanOrg(ctx context.Context, orgID int, banned bool) error {
	err := s.repo.BanOrg(ctx, orgID, banned)
	if err != nil {
		s.log.Error("failed to ban/unban organisation", "error", err, "org_id", orgID, "banned", banned)
		return err
	}
	return nil
}

func (s *service) ChangeBalance(ctx context.Context, orgID int, delta float64) error {
	return s.repo.ChangeBalance(ctx, orgID, delta)
}
