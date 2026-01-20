package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Starostina-elena/investment_platform/services/transactions/clients"
)

type Transaction struct {
	ID        int                `json:"id"`
	FromType  clients.EntityType `json:"from_type"`
	FromID    int                `json:"from_id"`
	ToType    clients.EntityType `json:"to_type"`
	ToID      int                `json:"to_id"`
	Amount    float64            `json:"amount"`
	CreatedAt time.Time          `json:"created_at"`
}

type Repo interface {
	Create(ctx context.Context, t *Transaction) (int, error)
}

type Service interface {
	Transfer(ctx context.Context, fromType, toType clients.EntityType, fromID, toID int, amount float64) (*Transaction, error)
}

type service struct {
	repo    Repo
	clients *clients.BalanceClient
	log     slog.Logger
}

func NewService(repo Repo, bc *clients.BalanceClient, log slog.Logger) Service {
	return &service{repo: repo, clients: bc, log: log}
}

func (s *service) Transfer(ctx context.Context, fromType, toType clients.EntityType, fromID, toID int, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	s.log.Info("starting transfer", "from", fromType, "from_id", fromID, "to", toType, "to_id", toID, "amount", amount)

	// --- 1. Списание средств (Deduct) ---
	// Если здесь ошибка (недостаточно средств), процесс прерывается
	err := s.clients.ChangeBalance(ctx, fromType, fromID, -amount)
	if err != nil {
		s.log.Error("failed to deduct funds", "error", err)
		return nil, fmt.Errorf("transaction failed: %v", err)
	}

	// --- 2. Начисление средств (Add) ---
	err = s.clients.ChangeBalance(ctx, toType, toID, amount)
	if err != nil {
		s.log.Error("failed to add funds, starting rollback", "error", err)

		// --- COMPENSATING TRANSACTION (ROLLBACK) ---
		// Возвращаем деньги отправителю
		rbErr := s.clients.ChangeBalance(ctx, fromType, fromID, amount)
		if rbErr != nil {
			// Это критическая ситуация, требует ручного вмешательства администратора
			s.log.Error("CRITICAL: ROLLBACK FAILED", "from_type", fromType, "from_id", fromID, "amount", amount, "error", rbErr)
			return nil, fmt.Errorf("system error: money stuck, contact support")
		}

		return nil, fmt.Errorf("transaction failed at destination: %v", err)
	}

	// --- 3. Сохранение истории ---
	t := &Transaction{
		FromType:  fromType,
		FromID:    fromID,
		ToType:    toType,
		ToID:      toID,
		Amount:    amount,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, t)
	if err != nil {
		// Деньги переведены, но история не сохранилась. Не критично для балансов, но плохо для отчетности.
		s.log.Error("transaction successful but failed to save record", "error", err)
	}
	t.ID = id

	return t, nil
}
