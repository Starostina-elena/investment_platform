package yookassa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type Client struct {
	shopID    string
	secretKey string
	baseURL   string
	client    *http.Client
}

func NewClient() *Client {
	shopID := os.Getenv("YOOKASSA_SHOP_ID")
	secretKey := os.Getenv("YOOKASSA_SECRET_KEY")

	// Логируем для отладки (без показа полного ключа)
	if shopID == "" || secretKey == "" {
		fmt.Println("WARNING: YOOKASSA credentials are empty!")
	} else {
		fmt.Printf("YooKassa client initialized:\n")
		fmt.Printf("  shopID: '%s' (len=%d)\n", shopID, len(shopID))
		fmt.Printf("  secretKey: '%s...' (len=%d)\n", secretKey[:min(15, len(secretKey))], len(secretKey))

		// Проверка на подозрительные символы
		if len(secretKey) > 5 && secretKey[5] == '*' {
			fmt.Println("  ⚠️  WARNING: Secret key has '*' at position 5 - this looks suspicious!")
			fmt.Println("  ⚠️  Please check if you copied the key correctly from YooKassa dashboard")
		}
	}

	return &Client{
		shopID:    shopID,
		secretKey: secretKey,
		baseURL:   "https://api.yookassa.ru/v3",
		client:    &http.Client{},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PaymentRequest - запрос на создание платежа (пополнение баланса)
type PaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type      string `json:"type"` // "redirect"
		ReturnURL string `json:"return_url,omitempty"`
	} `json:"confirmation"`
	Capture     bool   `json:"capture"` // true - деньги сразу списываются
	Description string `json:"description"`
	Metadata    struct {
		UserID int `json:"user_id"`
	} `json:"metadata,omitempty"`
}

// PaymentResponse - ответ при создании платежа
type PaymentResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"` // pending, waiting_for_capture, succeeded, canceled
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Description  string `json:"description"`
	Confirmation struct {
		Type            string `json:"type"`
		ConfirmationURL string `json:"confirmation_url"` // Ссылка для оплаты
	} `json:"confirmation"`
	CreatedAt string `json:"created_at"`
	Paid      bool   `json:"paid"`
	Metadata  struct {
		UserID string `json:"user_id"` // YooKassa возвращает как строку
	} `json:"metadata"`
}

// CreatePayment создает платеж и возвращает ссылку для оплаты
func (c *Client) CreatePayment(ctx context.Context, amount float64, userID int, description string, returnURL string) (*PaymentResponse, error) {
	req := PaymentRequest{
		Capture:     true,
		Description: description,
	}
	req.Amount.Value = fmt.Sprintf("%.2f", amount)
	req.Amount.Currency = "RUB"
	req.Confirmation.Type = "redirect"
	req.Confirmation.ReturnURL = returnURL
	req.Metadata.UserID = userID

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/payments", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Basic Auth: shopId в качестве username, secretKey в качестве password
	// Для тестового режима secretKey начинается с "test_"
	httpReq.SetBasicAuth(c.shopID, c.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")
	// Генерируем UUID для Idempotence-Key как в SDK
	idempotenceKey := uuid.New().String()
	httpReq.Header.Set("Idempotence-Key", idempotenceKey)

	// Логируем запрос для отладки
	authHeader := httpReq.Header.Get("Authorization")
	fmt.Printf("YooKassa CreatePayment: POST %s/payments\n", c.baseURL)
	fmt.Printf("  shopID: %s\n", c.shopID)
	fmt.Printf("  secretKey prefix: %s\n", c.secretKey[:min(5, len(c.secretKey))])
	fmt.Printf("  Idempotence-Key: %s\n", idempotenceKey)
	fmt.Printf("  Authorization header exists: %v\n", authHeader != "")
	fmt.Printf("  Content-Type: %s\n", httpReq.Header.Get("Content-Type"))
	fmt.Printf("  Request body: %s\n", string(body))

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YooKassa API error: %s, body: %s", resp.Status, string(respBody))
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(respBody, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &paymentResp, nil
}

// GetPayment получает информацию о платеже
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*PaymentResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/payments/"+paymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.SetBasicAuth(c.shopID, c.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YooKassa API error: %s, body: %s", resp.Status, string(respBody))
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(respBody, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &paymentResp, nil
}
