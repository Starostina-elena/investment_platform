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

type ProjectData struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	MonetizationType       string  `json:"monetization_type"`
	Percent                float64 `json:"percent"`
	CurrentMoney           float64 `json:"current_money"`
	WantedMoney            float64 `json:"wanted_money"`
	MoneyRequiredToPayback float64 `json:"money_required_to_payback"`
	CreatorID              int     `json:"creator_id"`
	IsCompleted            bool    `json:"is_completed"`
}

func NewProjectClient(log slog.Logger) *ProjectClient {
	url := os.Getenv("PROJECT_SERVICE_URL")
	if url == "" {
		url = "http://project:8104"
	}
	return &ProjectClient{
		url:    url,
		client: &http.Client{},
		log:    log,
	}
}

func (pc *ProjectClient) GetProject(ctx context.Context, projectID int) (*ProjectData, error) {
	url := fmt.Sprintf("%s/%d", pc.url, projectID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	if err != nil {
		pc.log.Error("failed to get project", "error", err, "project_id", projectID)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		pc.log.Error("project service error", "status", resp.StatusCode, "project_id", projectID)
		return nil, fmt.Errorf("project service error: status %d", resp.StatusCode)
	}

	var project ProjectData
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		pc.log.Error("failed to decode project response", "error", err)
		return nil, err
	}

	return &project, nil
}

func (pc *ProjectClient) UpdateMoneyRequiredToPayback(ctx context.Context, projectID int, amount float64) error {
	url := fmt.Sprintf("%s/%d/money-required-payback", pc.url, projectID)

	reqBody, _ := json.Marshal(map[string]interface{}{"amount": amount})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.client.Do(req)
	if err != nil {
		pc.log.Error("failed to update money required to payback", "error", err, "project_id", projectID)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		pc.log.Error("project service error", "status", resp.StatusCode, "project_id", projectID)
		return fmt.Errorf("project service error: status %d", resp.StatusCode)
	}

	return nil
}

func (pc *ProjectClient) GetProjectOwnerEmail(ctx context.Context, projectID int) (string, error) {
	project, err := pc.GetProject(ctx, projectID)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/internal/user/%d/email", os.Getenv("USER_SERVICE_URL"), project.CreatorID)
	if userURL := os.Getenv("USER_SERVICE_URL"); userURL == "" {
		url = fmt.Sprintf("http://user:8101/internal/user/%d/email", project.CreatorID)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := pc.client.Do(req)
	if err != nil {
		pc.log.Error("failed to get user email", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		pc.log.Error("user service error", "status", resp.StatusCode)
		return "", fmt.Errorf("user service error: status %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		pc.log.Error("failed to decode user response", "error", err)
		return "", err
	}

	return result["email"], nil
}
