package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TransactionClient struct {
	url    string
	client *http.Client
}

func NewTransactionClient() *TransactionClient {
	url := os.Getenv("TRANSACTION_SERVICE_URL")
	if url == "" {
		url = "http://transactions:8103"
	}
	return &TransactionClient{
		url:    url,
		client: &http.Client{},
	}
}

func (tc *TransactionClient) Deposit(ctx context.Context, toType string, toID int, amount float64) error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"from_type": "external", // Важно! Это помечает ввод средств
		"from_id":   0,          // ID системы/шлюза
		"to_type":   toType,     // "user" или "org"
		"to_id":     toID,
		"amount":    amount,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", tc.url+"/transfer", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := tc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("transaction service error: %d", resp.StatusCode)
	}
	return nil
}
