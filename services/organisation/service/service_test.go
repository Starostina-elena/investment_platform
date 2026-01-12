package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
)

type mockService struct {
	banOrgFunc       func(ctx context.Context, orgID int, banned bool) error
	getUsersOrgsFunc func(ctx context.Context, userID int) ([]core.Org, error)
}

func (m *mockService) Create(ctx context.Context, org core.Org) (*core.Org, error) {
	return &org, nil
}

func (m *mockService) Get(ctx context.Context, id int) (*core.Org, error) {
	return nil, nil
}

func (m *mockService) Update(ctx context.Context, org core.Org, userRequestedId int) (*core.Org, error) {
	return &org, nil
}

func (m *mockService) UploadAvatar(ctx context.Context, orgID int, userID int, file interface{}, fileHeader interface{}) (string, error) {
	return "", nil
}

func (m *mockService) DeleteAvatar(ctx context.Context, orgID int, userID int, avatarPath string) error {
	return nil
}

func (m *mockService) UpdateAvatarPath(ctx context.Context, orgID int, avatarPath string) error {
	return nil
}

func (m *mockService) UploadDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType, file interface{}, fileHeader interface{}) (string, error) {
	return "", nil
}

func (m *mockService) DeleteDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType) error {
	return nil
}

func (m *mockService) DownloadDoc(ctx context.Context, orgID int, userID int, isAdmin bool, docType core.OrgDocType) ([]byte, string, error) {
	return nil, "", nil
}

func (m *mockService) BanOrg(ctx context.Context, orgID int, banned bool) error {
	if m.banOrgFunc != nil {
		return m.banOrgFunc(ctx, orgID, banned)
	}
	return nil
}

func (m *mockService) GetUsersOrgs(ctx context.Context, userID int) ([]core.Org, error) {
	if m.getUsersOrgsFunc != nil {
		return m.getUsersOrgsFunc(ctx, userID)
	}
	return nil, nil
}

func TestBanOrg_Success(t *testing.T) {
	ms := &mockService{
		banOrgFunc: func(ctx context.Context, orgID int, banned bool) error {
			return nil
		},
	}

	err := ms.BanOrg(context.Background(), 1, true)
	if err != nil {
		t.Errorf("BanOrg() error = %v", err)
	}
}

func TestBanOrg_Error(t *testing.T) {
	ms := &mockService{
		banOrgFunc: func(ctx context.Context, orgID int, banned bool) error {
			return errors.New("database error")
		},
	}

	err := ms.BanOrg(context.Background(), 1, true)
	if err == nil {
		t.Error("BanOrg() expected error")
	}
}

func TestGetUsersOrgs_Success(t *testing.T) {
	ms := &mockService{
		getUsersOrgsFunc: func(ctx context.Context, userID int) ([]core.Org, error) {
			return []core.Org{
				{
					OrgBase: core.OrgBase{
						ID:      1,
						Name:    "Test Org",
						OrgType: core.OrgTypePhys,
					},
				},
			}, nil
		},
	}

	orgs, err := ms.GetUsersOrgs(context.Background(), 1)
	if err != nil {
		t.Errorf("GetUsersOrgs() error = %v", err)
	}
	if len(orgs) != 1 {
		t.Errorf("GetUsersOrgs() expected 1 org, got %d", len(orgs))
	}
}

func TestGetUsersOrgs_Empty(t *testing.T) {
	ms := &mockService{
		getUsersOrgsFunc: func(ctx context.Context, userID int) ([]core.Org, error) {
			return []core.Org{}, nil
		},
	}

	orgs, err := ms.GetUsersOrgs(context.Background(), 1)
	if err != nil {
		t.Errorf("GetUsersOrgs() error = %v", err)
	}
	if len(orgs) != 0 {
		t.Errorf("GetUsersOrgs() expected 0 orgs, got %d", len(orgs))
	}
}
