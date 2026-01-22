package yookassa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type Client struct {
	ShopID       string
	SecretKey    string
	AgentID      string
	PayoutAPIKey string
	APIURL       string
	PayoutAPIURL string
	HTTP         *http.Client
}

func NewClient() *Client {
	shopID := os.Getenv("YOOKASSA_SHOP_ID")
	secretKey := os.Getenv("YOOKASSA_SECRET_KEY")
	agentID := os.Getenv("YOOKASSA_AGENT_ID")
	payoutAPIKey := os.Getenv("YOOKASSA_PAYOUT_API_KEY")

	fmt.Printf("YooKassa Client initialized:\n")
	fmt.Printf("  Shop ID: %s\n", shopID)
	fmt.Printf("  Agent ID: %s\n", agentID)

	return &Client{
		ShopID:       shopID,
		SecretKey:    secretKey,
		AgentID:      agentID,
		PayoutAPIKey: payoutAPIKey,
		APIURL:       "https://api.yookassa.ru/v3/payments",
		PayoutAPIURL: "https://api.yookassa.ru/v3/payouts",
		HTTP:         &http.Client{},
	}
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type CreatePaymentRequest struct {
	Amount       Amount       `json:"amount"`
	Capture      bool         `json:"capture"`
	Confirmation Confirmation `json:"confirmation"`
	Description  string       `json:"description"`
}

type CreatePaymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Paid         bool   `json:"paid"`
	Confirmation struct {
		Type            string `json:"type"`
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}

func (c *Client) CreatePayment(amount string, description string, returnURL string) (*CreatePaymentResponse, error) {
	reqBody := CreatePaymentRequest{
		Amount: Amount{
			Value:    amount,
			Currency: "RUB",
		},
		Capture: true,
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnURL: returnURL,
		},
		Description: description,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", c.APIURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.ShopID + ":" + c.SecretKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Idempotence-Key", uuid.New().String())
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("YooKassa CreatePayment:\n")
	fmt.Printf("  Auth header: Basic %s\n", auth[:min(20, len(auth))]+"...")
	fmt.Printf("  URL: %s\n", c.APIURL)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yookassa error: %s", string(respBody))
	}

	var result CreatePaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetPayment(paymentID string) (*CreatePaymentResponse, error) {
	req, err := http.NewRequest("GET", c.APIURL+"/"+paymentID, nil)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.ShopID + ":" + c.SecretKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yookassa error: %s", string(respBody))
	}

	var result CreatePaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type CreatePayoutRequest struct {
	Amount      Amount            `json:"amount"`
	Description string            `json:"description"`
	PayoutToken string            `json:"payout_token"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type PayoutResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (c *Client) CreatePayout(amount string, description string, payoutToken string) (*PayoutResponse, error) {
	reqBody := CreatePayoutRequest{
		Amount: Amount{
			Value:    amount,
			Currency: "RUB",
		},
		Description: description,
		PayoutToken: payoutToken,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	fmt.Printf("YooKassa CreatePayout request JSON:\n%s\n", string(bodyBytes))
	req, err := http.NewRequest("POST", c.PayoutAPIURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.AgentID + ":" + c.PayoutAPIKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Idempotence-Key", uuid.New().String())
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yookassa payout error: %s", string(respBody))
	}

	var result PayoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetPayout(payoutID string) (*PayoutResponse, error) {
	payoutURL := "https://api.yookassa.ru/v3/payouts/" + payoutID

	req, err := http.NewRequest("GET", payoutURL, nil)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.AgentID + ":" + c.PayoutAPIKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yookassa payout error: %s", string(respBody))
	}

	var result PayoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
