package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type EntityType string

const (
	TypeUser    EntityType = "user"
	TypeOrg     EntityType = "org"
	TypeProject EntityType = "project"
)

type BalanceClient struct {
	userUrl    string
	orgUrl     string
	projectUrl string
	client     *http.Client
	log        slog.Logger
}

func NewBalanceClient(log slog.Logger) *BalanceClient {
	return &BalanceClient{
		userUrl:    getEnv("USER_SERVICE_URL", "http://user:8101"),
		orgUrl:     getEnv("ORG_SERVICE_URL", "http://organisation:8102"),
		projectUrl: getEnv("PROJECT_SERVICE_URL", "http://project:8104"),
		client:     &http.Client{},
		log:        log,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (bc *BalanceClient) ChangeBalance(ctx context.Context, entityType EntityType, id int, delta float64) error {
	var url string
	var reqBody []byte

	switch entityType {
	case TypeUser:
		url = fmt.Sprintf("%s/internal/balance", bc.userUrl)
		reqBody, _ = json.Marshal(map[string]interface{}{"id": id, "delta": delta})
	case TypeOrg:
		url = fmt.Sprintf("%s/internal/balance", bc.orgUrl)
		reqBody, _ = json.Marshal(map[string]interface{}{"id": id, "delta": delta})
	case TypeProject:
		url = fmt.Sprintf("%s/%d/funds", bc.projectUrl, id)
		reqBody, _ = json.Marshal(map[string]interface{}{"amount": delta})
	default:
		return fmt.Errorf("unknown entity type: %s", entityType)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := bc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bc.log.Error("remote service error", "url", url, "status", resp.StatusCode)
		return fmt.Errorf("balance update failed for %s id=%d: status %d", entityType, id, resp.StatusCode)
	}

	return nil
}
