package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Starostina-elena/investment_platform/services/payment/clients"
	"github.com/Starostina-elena/investment_platform/services/payment/core"
	"github.com/Starostina-elena/investment_platform/services/payment/repo"
	"github.com/Starostina-elena/investment_platform/services/payment/yookassa"
	"github.com/google/uuid"
)

type Service struct {
	repo     *repo.Repo
	yookassa *yookassa.Client
	txClient *clients.TransactionClient
	log      slog.Logger
}

func NewService(repo *repo.Repo, yc *yookassa.Client, tc *clients.TransactionClient, log slog.Logger) *Service {
	return &Service{repo: repo, yookassa: yc, txClient: tc, log: log}
}

// InitPayment создает платеж в ЮКассе и сохраняет в БД
func (s *Service) InitPayment(ctx context.Context, entityType string, entityID int, amount float64, returnURL string) (string, error) {
	amountStr := fmt.Sprintf("%.2f", amount)
	desc := fmt.Sprintf("Пополнение кошелька %s #%d", entityType, entityID)

	// 1. Запрос в ЮКассу
	yooResp, err := s.yookassa.CreatePayment(amountStr, desc, returnURL)
	if err != nil {
		s.log.Error("yookassa create failed", "error", err)
		return "", err
	}

	// 2. Сохранение в БД
	payment := &core.Payment{
		ID:         uuid.New().String(),
		ExternalID: yooResp.ID,
		Amount:     amount,
		EntityID:   entityID,
		EntityType: entityType,
		Status:     core.StatusPending,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		s.log.Error("db save failed", "error", err)
		return "", err
	}

	return yooResp.Confirmation.ConfirmationURL, nil
}

// ProcessWebhook обрабатывает уведомление от ЮКассы
func (s *Service) ProcessWebhook(ctx context.Context, eventType string, object map[string]interface{}) error {
	if eventType != "payment.succeeded" {
		return nil // Нас интересуют только успешные оплаты
	}

	externalID, _ := object["id"].(string)
	if externalID == "" {
		return fmt.Errorf("empty payment id in webhook")
	}

	// 1. Получаем платеж из БД
	payment, err := s.repo.GetByExternalID(ctx, externalID)
	if err != nil {
		s.log.Error("payment not found", "id", externalID)
		return err
	}

	// Идемпотентность: если уже успешно, выходим
	if payment.Status == core.StatusSucceeded {
		return nil
	}

	// 2. Обновляем статус в БД
	if err := s.repo.UpdateStatus(ctx, externalID, core.StatusSucceeded); err != nil {
		return err
	}

	// 3. Начисляем деньги через Transaction Service
	s.log.Info("crediting wallet", "entity_type", payment.EntityType, "entity_id", payment.EntityID, "amount", payment.Amount)
	if err := s.txClient.Deposit(ctx, payment.EntityType, payment.EntityID, payment.Amount); err != nil {
		s.log.Error("CRITICAL: failed to deposit money after success payment", "error", err, "payment_id", payment.ID)
		return err
	}

	return nil
}

// CheckPayment проверяет статус платежа в ЮКассе и обновляет БД
func (s *Service) CheckPayment(ctx context.Context, paymentID string) (*core.Payment, error) {
	// 1. Получаем платеж из БД по ID платежа сервиса
	payment, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		s.log.Error("payment not found", "id", paymentID)
		return nil, fmt.Errorf("payment not found")
	}

	// 2. Проверяем статус в ЮКассе
	yooPayment, err := s.yookassa.GetPayment(payment.ExternalID)
	if err != nil {
		s.log.Error("failed to get payment from yookassa", "error", err, "external_id", payment.ExternalID)
		return nil, err
	}

	// 3. Если статус изменился на succeeded, обновляем БД и начисляем деньги
	if yooPayment.Status == "succeeded" && payment.Status != core.StatusSucceeded {
		if err := s.repo.UpdateStatus(ctx, payment.ExternalID, core.StatusSucceeded); err != nil {
			s.log.Error("failed to update payment status", "error", err)
			return nil, err
		}

		// Начисляем деньги
		s.log.Info("crediting wallet from check", "entity_type", payment.EntityType, "entity_id", payment.EntityID, "amount", payment.Amount)
		if err := s.txClient.Deposit(ctx, payment.EntityType, payment.EntityID, payment.Amount); err != nil {
			s.log.Error("CRITICAL: failed to deposit money after check", "error", err, "payment_id", payment.ID)
			return nil, err
		}

		payment.Status = core.StatusSucceeded
	}

	return payment, nil
}

// ProcessPendingPayments проверяет все pending платежи в YooKassa
func (s *Service) ProcessPendingPayments(ctx context.Context) error {
	payments, err := s.repo.GetPendingPayments(ctx)
	if err != nil {
		s.log.Error("failed to get pending payments", "error", err)
		return err
	}

	for _, payment := range payments {
		yooPayment, err := s.yookassa.GetPayment(payment.ExternalID)
		if err != nil {
			s.log.Error("failed to check payment status", "error", err, "payment_id", payment.ID)
			continue
		}

		if yooPayment.Status == "succeeded" {
			s.log.Info("payment succeeded, crediting wallet", "payment_id", payment.ID, "external_id", payment.ExternalID)

			if err := s.repo.UpdateStatus(ctx, payment.ExternalID, core.StatusSucceeded); err != nil {
				s.log.Error("failed to update status", "error", err)
				continue
			}

			if err := s.txClient.Deposit(ctx, payment.EntityType, payment.EntityID, payment.Amount); err != nil {
				s.log.Error("failed to deposit money", "error", err, "payment_id", payment.ID)
				continue
			}
		}
	}

	return nil
}
