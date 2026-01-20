package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type NotificationClient struct {
	url    string
	client *http.Client
	log    slog.Logger
}

func NewNotificationClient(log slog.Logger) *NotificationClient {
	url := os.Getenv("NOTIFICATION_SERVICE_URL")
	if url == "" {
		url = "http://notification:8083"
	}
	return &NotificationClient{
		url:    url,
		client: &http.Client{},
		log:    log,
	}
}

func (nc *NotificationClient) SendEmail(email, notifType, projectName string, amount float64) error {
	if notifType == "project_goal_reached" {
		notifType = "project_closed"
	}

	payload := map[string]interface{}{
		"email":        email,
		"type":         notifType,
		"project_name": projectName,
		"amount":       amount,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(nc.url+"/send", "application/json", bytes.NewReader(body))
	if err != nil {
		nc.log.Error("failed to send email", "error", err, "to", email)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		nc.log.Error("email service returned error", "status", resp.StatusCode, "to", email)
		return fmt.Errorf("notification service error: status %d", resp.StatusCode)
	}

	return nil
}
