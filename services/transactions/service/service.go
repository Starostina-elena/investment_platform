package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Starostina-elena/investment_platform/services/transactions/clients"
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

// Добавляем отсутствующий интерфейс Repo
type Repo interface {
	Create(ctx context.Context, t *Transaction) (int, error)
}

type Service interface {
	Invest(ctx context.Context, userID, projectID int, amount float64, method string) (*Transaction, error)
}

type service struct {
	repo          Repo
	projectClient *clients.ProjectClient
	log           slog.Logger
}

func NewService(repo Repo, pc *clients.ProjectClient, log slog.Logger) Service {
	return &service{repo: repo, projectClient: pc, log: log}
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
