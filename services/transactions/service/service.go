package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Starostina-elena/investment_platform/services/transactions/clients"
	"github.com/Starostina-elena/investment_platform/services/transactions/clients/yookassa"
)

type Transaction struct {
	ID        int       `json:"id"`
	FromID    int       `json:"from_id"`
	ToID      int       `json:"to_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	Method    string    `json:"method"` // sbp, yookassa
	CreatedAt time.Time `json:"created_at"`
}

type Repo interface {
	Create(ctx context.Context, t *Transaction) (int, error)
}

type Service interface {
	Invest(ctx context.Context, userID, projectID int, amount float64, method string) (*Transaction, error)
	CreateDeposit(ctx context.Context, userID int, amount float64, returnURL string) (string, string, error) // возвращает payment_id и confirmation_url
	CheckDeposit(ctx context.Context, paymentID string) (*Transaction, error)
}

type service struct {
	repo          Repo
	projectClient *clients.ProjectClient
	yookassa      *yookassa.Client
	log           slog.Logger
}

func NewService(repo Repo, pc *clients.ProjectClient, log slog.Logger) Service {
	return &service{
		repo:          repo,
		projectClient: pc,
		yookassa:      yookassa.NewClient(),
		log:           log,
	}
}

func (s *service) Invest(ctx context.Context, userID, projectID int, amount float64, method string) (*Transaction, error) {
	s.log.Info("processing payment", "method", method, "amount", amount)

	// Имитация задержки
	time.Sleep(500 * time.Millisecond)

	if method != "sbp" && method != "yookassa" {
		return nil, fmt.Errorf("unknown payment method: %s", method)
	}

	t := &Transaction{
		FromID:    userID,
		ToID:      projectID,
		Amount:    amount,
		Type:      "user_to_project",
		Method:    method,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, t)
	if err != nil {
		s.log.Error("failed to create transaction record", "error", err)
		return nil, err
	}
	t.ID = id

	// Вызов сервиса проектов для начисления средств
	if err := s.projectClient.AddFunds(ctx, projectID, amount); err != nil {
		s.log.Error("CRITICAL: payment succeeded but project balance update failed", "tx_id", id, "error", err)
		// В реальном мире здесь нужно делать возврат средств или ретрай
		return nil, fmt.Errorf("funds transfer failed")
	}

	return t, nil
}

// CreateDeposit создает платеж для пополнения баланса и возвращает ссылку для оплаты
func (s *service) CreateDeposit(ctx context.Context, userID int, amount float64, returnURL string) (string, string, error) {
	s.log.Info("creating deposit payment", "user_id", userID, "amount", amount)

	if amount <= 0 {
		return "", "", fmt.Errorf("amount must be positive")
	}

	payment, err := s.yookassa.CreatePayment(ctx, amount, userID, fmt.Sprintf("Пополнение баланса пользователя %d", userID), returnURL)
	if err != nil {
		s.log.Error("failed to create YooKassa payment", "error", err, "user_id", userID)
		return "", "", fmt.Errorf("failed to create payment: %w", err)
	}

	s.log.Info("YooKassa payment created", "payment_id", payment.ID, "status", payment.Status, "confirmation_url", payment.Confirmation.ConfirmationURL)

	// Пока не сохраняем транзакцию - она будет создана после успешной оплаты в CheckDeposit
	return payment.ID, payment.Confirmation.ConfirmationURL, nil
}

// CheckDeposit проверяет статус платежа и создает транзакцию если оплата прошла
func (s *service) CheckDeposit(ctx context.Context, paymentID string) (*Transaction, error) {
	s.log.Info("checking deposit payment status", "payment_id", paymentID)

	payment, err := s.yookassa.GetPayment(ctx, paymentID)
	if err != nil {
		s.log.Error("failed to get payment status", "error", err, "payment_id", paymentID)
		return nil, fmt.Errorf("failed to check payment: %w", err)
	}

	s.log.Info("payment status retrieved", "payment_id", paymentID, "status", payment.Status, "paid", payment.Paid)

	if payment.Status != "succeeded" || !payment.Paid {
		return nil, fmt.Errorf("payment not completed yet, status: %s", payment.Status)
	}

	// Парсим UserID из строки в int
	userID, err := strconv.Atoi(payment.Metadata.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user_id: %w", err)
	}

	// Парсим сумму
	var amountFloat float64
	if _, err := fmt.Sscanf(payment.Amount.Value, "%f", &amountFloat); err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	// Создаем транзакцию пополнения
	t := &Transaction{
		FromID:    0, // 0 означает внешнее пополнение
		ToID:      userID,
		Amount:    amountFloat,
		Type:      "user_deposit",
		Method:    "yookassa",
		CreatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, t)
	if err != nil {
		s.log.Error("failed to create deposit transaction", "error", err)
		return nil, err
	}
	t.ID = id

	s.log.Info("deposit processed successfully", "tx_id", id, "user_id", userID, "yookassa_id", paymentID)

	// TODO: Обновить баланс пользователя через user service

	return t, nil
}
