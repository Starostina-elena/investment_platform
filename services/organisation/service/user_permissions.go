package service

import (
	"context"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
)

func (s *service) CheckUserOrgPermission(ctx context.Context, orgID int, userID int, permission core.OrgPermission) (bool, error) {
	permissions, err := s.repo.GetUserOrgPermissions(ctx, orgID, userID)
	if err != nil {
		s.log.Error("failed to check user organisation permission", "org_id", orgID, "user_id", userID, "permission", permission, "error", err)
		return false, err
	}
	return permissions[string(permission)], nil
}

func (s *service) AddEmployee(ctx context.Context, orgID int, userRequested int, userID int, orgAccMgmt, moneyMgmt, projMgmt bool) error {
	allowed, err := s.CheckUserOrgPermission(ctx, orgID, userRequested, "org_account_management")
	if err != nil {
		return err
	}
	if !allowed {
		return core.ErrNotAuthorized
	}
	return s.repo.AddEmployee(ctx, orgID, userID, orgAccMgmt, moneyMgmt, projMgmt)
}

func (s *service) GetOrgEmployees(ctx context.Context, orgID int) ([]core.OrgEmployee, error) {
	return s.repo.GetEmployees(ctx, orgID)
}

func (s *service) UpdateEmployeePermissions(ctx context.Context, orgID int, userRequested int, userID int, orgAccMgmt, moneyMgmt, projMgmt bool) error {
	allowed, err := s.CheckUserOrgPermission(ctx, orgID, userRequested, "org_account_management")
	if err != nil {
		return err
	}
	if !allowed {
		return core.ErrNotAuthorized
	}
	return s.repo.UpdateEmployeePermissions(ctx, orgID, userID, orgAccMgmt, moneyMgmt, projMgmt)
}

func (s *service) DeleteEmployee(ctx context.Context, orgID int, userRequested int, userID int) error {
	allowed, err := s.CheckUserOrgPermission(ctx, orgID, userRequested, "org_account_management")
	if err != nil {
		return err
	}
	if !allowed {
		return core.ErrNotAuthorized
	}
	return s.repo.DeleteEmployee(ctx, orgID, userID)
}
