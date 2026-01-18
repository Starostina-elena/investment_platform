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
