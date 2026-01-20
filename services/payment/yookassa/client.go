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
	ShopID    string
	SecretKey string
	APIURL    string
	HTTP      *http.Client
}

func NewClient() *Client {
	return &Client{
		ShopID:    os.Getenv("YOOKASSA_SHOP_ID"),
		SecretKey: os.Getenv("YOOKASSA_SECRET_KEY"),
		APIURL:    "https://api.yookassa.ru/v3/payments",
		HTTP:      &http.Client{},
	}
}

// Структуры для запросов/ответов ЮKassa
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

	// Auth headers
	auth := base64.StdEncoding.EncodeToString([]byte(c.ShopID + ":" + c.SecretKey))
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
		return nil, fmt.Errorf("yookassa error: %s", string(respBody))
	}

	var result CreatePaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
