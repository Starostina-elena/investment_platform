package service

import (
	"context"
	"fmt"

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

func (s *service) TransferOwnership(ctx context.Context, orgID int, userRequested int, newOwnerID int) error {
	org, err := s.Get(ctx, orgID)
	if err != nil {
		return err
	}
	if org.OwnerId != userRequested {
		s.log.Warn("non-owner tried to transfer ownership", "org_id", orgID, "user_id", userRequested, "owner_id", org.OwnerId)
		return core.ErrNotAuthorized
	}

	if org.IsBanned {
		s.log.Warn("attempt to transfer ownership of banned organisation", "org_id", orgID)
		return core.ErrOrgBanned
	}

	if org.OwnerId == newOwnerID {
		s.log.Warn("attempt to transfer ownership to current owner", "org_id", orgID, "owner_id", org.OwnerId)
		return fmt.Errorf("new owner is the same as current owner")
	}

	s.log.Info("transferring ownership", "org_id", orgID, "old_owner", org.OwnerId, "new_owner", newOwnerID)

	return s.repo.TransferOwnership(ctx, orgID, org.OwnerId, newOwnerID)
}
