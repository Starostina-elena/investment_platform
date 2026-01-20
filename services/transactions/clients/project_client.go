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

type ProjectClient struct {
	url    string
	client *http.Client
	log    slog.Logger
}

func NewProjectClient(log slog.Logger) *ProjectClient {
	url := os.Getenv("PROJECT_SERVICE_URL")
	if url == "" {
		url = "http://project:8104" // дефолт для docker-compose
	}
	return &ProjectClient{
		url:    url,
		client: &http.Client{},
		log:    log,
	}
}

func (pc *ProjectClient) AddFunds(ctx context.Context, projectID int, amount float64) error {
	reqBody, _ := json.Marshal(map[string]float64{
		"amount": amount,
	})

	url := fmt.Sprintf("%s/%d/funds", pc.url, projectID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		pc.log.Error("project service returned error", "status", resp.StatusCode)
		return fmt.Errorf("project service error: %d", resp.StatusCode)
	}

	return nil
}
