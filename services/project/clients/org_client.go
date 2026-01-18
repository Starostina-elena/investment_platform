package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type OrgClient struct {
	url    string
	client *http.Client
	log    slog.Logger
}

func NewOrgClient(url string, log slog.Logger) *OrgClient {
	return &OrgClient{
		url:    url,
		client: &http.Client{},
		log:    log,
	}
}

func (oc *OrgClient) CheckUserOrgPermission(ctx context.Context, orgID int, userID int, permission string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%d/rights/%d/%s", oc.url, orgID, userID, permission), nil)
	if err != nil {
		return false, err
	}

	resp, err := oc.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("failed to get organisation permissions")
	}

	var perms map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&perms); err != nil {
		return false, err
	}

	allowed, exists := perms["allowed"]
	if !exists {
		return false, errors.New("incorrect response format")
	}

	return allowed, nil
}
