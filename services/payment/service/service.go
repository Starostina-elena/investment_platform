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

func (s *Service) InitPayment(ctx context.Context, entityType string, entityID int, amount float64, returnURL string) (string, error) {
	amountStr := fmt.Sprintf("%.2f", amount)
	desc := fmt.Sprintf("Пополнение кошелька %s #%d", entityType, entityID)

	yooResp, err := s.yookassa.CreatePayment(amountStr, desc, returnURL)
	if err != nil {
		s.log.Error("yookassa create failed", "error", err)
		return "", err
	}

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

func (s *Service) ProcessWebhook(ctx context.Context, eventType string, object map[string]interface{}) error {
	if eventType != "payment.succeeded" {
		return nil
	}

	externalID, _ := object["id"].(string)
	if externalID == "" {
		return fmt.Errorf("empty payment id in webhook")
	}

	payment, err := s.repo.GetByExternalID(ctx, externalID)
	if err != nil {
		s.log.Error("payment not found", "id", externalID)
		return err
	}

	if payment.Status == core.StatusSucceeded {
		return nil
	}

	if err := s.repo.UpdateStatus(ctx, externalID, core.StatusSucceeded); err != nil {
		return err
	}

	s.log.Info("crediting wallet", "entity_type", payment.EntityType, "entity_id", payment.EntityID, "amount", payment.Amount)
	if err := s.txClient.Deposit(ctx, payment.EntityType, payment.EntityID, payment.Amount); err != nil {
		s.log.Error("CRITICAL: failed to deposit money after success payment", "error", err, "payment_id", payment.ID)
		return err
	}

	return nil
}

func (s *Service) CheckPayment(ctx context.Context, paymentID string) (*core.Payment, error) {
	payment, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		s.log.Error("payment not found", "id", paymentID)
		return nil, fmt.Errorf("payment not found")
	}

	yooPayment, err := s.yookassa.GetPayment(payment.ExternalID)
	if err != nil {
		s.log.Error("failed to get payment from yookassa", "error", err, "external_id", payment.ExternalID)
		return nil, err
	}

	if yooPayment.Status == "succeeded" && payment.Status != core.StatusSucceeded {
		if err := s.repo.UpdateStatus(ctx, payment.ExternalID, core.StatusSucceeded); err != nil {
			s.log.Error("failed to update payment status", "error", err)
			return nil, err
		}

		s.log.Info("crediting wallet from check", "entity_type", payment.EntityType, "entity_id", payment.EntityID, "amount", payment.Amount)
		if err := s.txClient.Deposit(ctx, payment.EntityType, payment.EntityID, payment.Amount); err != nil {
			s.log.Error("CRITICAL: failed to deposit money after check", "error", err, "payment_id", payment.ID)
			return nil, err
		}

		payment.Status = core.StatusSucceeded
	}

	return payment, nil
}

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

func (s *Service) InitWithdrawal(ctx context.Context, entityType string, entityID int, amount float64, destination string) (string, error) {
	amountStr := fmt.Sprintf("%.2f", amount)
	desc := fmt.Sprintf("Вывод средств %s #%d", entityType, entityID)

	yooResp, err := s.yookassa.CreatePayout(amountStr, desc, destination)
	if err != nil {
		s.log.Error("yookassa payout creation failed", "error", err)
		return "", err
	}

	withdrawal := &core.Withdrawal{
		ID:         uuid.New().String(),
		ExternalID: yooResp.ID,
		Amount:     amount,
		EntityID:   entityID,
		EntityType: entityType,
		Status:     core.WithdrawalPending,
	}

	if err := s.repo.CreateWithdrawal(ctx, withdrawal); err != nil {
		s.log.Error("db save withdrawal failed", "error", err)
		return "", err
	}

	if err := s.txClient.Withdraw(ctx, entityType, entityID, amount); err != nil {
		s.log.Error("failed to withdraw funds", "error", err)
		return "", err
	}

	return withdrawal.ID, nil
}

func (s *Service) ProcessWithdrawalWebhook(ctx context.Context, eventType string, object map[string]interface{}) error {
	if eventType != "payout.succeeded" && eventType != "payout.failed" {
		return nil
	}

	externalID, _ := object["id"].(string)
	if externalID == "" {
		return fmt.Errorf("empty payout id in webhook")
	}

	withdrawal, err := s.repo.GetWithdrawalByExternalID(ctx, externalID)
	if err != nil {
		s.log.Error("withdrawal not found", "id", externalID)
		return err
	}

	var status core.WithdrawalStatus
	if eventType == "payout.succeeded" {
		status = core.WithdrawalSucceeded
		s.log.Info("payout succeeded", "withdrawal_id", withdrawal.ID)
	} else {
		status = core.WithdrawalFailed
		s.log.Warn("payout failed", "withdrawal_id", withdrawal.ID)

		if err := s.txClient.Deposit(ctx, withdrawal.EntityType, withdrawal.EntityID, withdrawal.Amount); err != nil {
			s.log.Error("CRITICAL: failed to refund after failed payout", "error", err, "withdrawal_id", withdrawal.ID)
			return err
		}
	}

	if err := s.repo.UpdateWithdrawalStatus(ctx, externalID, status); err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckWithdrawal(ctx context.Context, withdrawalID string) (*core.Withdrawal, error) {
	withdrawal, err := s.repo.GetWithdrawalByID(ctx, withdrawalID)
	if err != nil {
		s.log.Error("withdrawal not found", "id", withdrawalID)
		return nil, fmt.Errorf("withdrawal not found")
	}

	yooPayout, err := s.yookassa.GetPayout(withdrawal.ExternalID)
	if err != nil {
		s.log.Error("failed to get payout from yookassa", "error", err, "external_id", withdrawal.ExternalID)
		return nil, err
	}

	if yooPayout.Status == "succeeded" && withdrawal.Status != core.WithdrawalSucceeded {
		if err := s.repo.UpdateWithdrawalStatus(ctx, withdrawal.ExternalID, core.WithdrawalSucceeded); err != nil {
			return nil, err
		}
		withdrawal.Status = core.WithdrawalSucceeded
	} else if yooPayout.Status == "failed" && withdrawal.Status != core.WithdrawalFailed {
		if err := s.repo.UpdateWithdrawalStatus(ctx, withdrawal.ExternalID, core.WithdrawalFailed); err != nil {
			return nil, err
		}

		if err := s.txClient.Deposit(ctx, withdrawal.EntityType, withdrawal.EntityID, withdrawal.Amount); err != nil {
			s.log.Error("failed to refund after failed payout", "error", err, "withdrawal_id", withdrawal.ID)
		}

		withdrawal.Status = core.WithdrawalFailed
	}

	return withdrawal, nil
}

func (s *Service) ProcessPendingWithdrawals(ctx context.Context) error {
	withdrawals, err := s.repo.GetPendingWithdrawals(ctx)
	if err != nil {
		s.log.Error("failed to get pending withdrawals", "error", err)
		return err
	}

	for _, withdrawal := range withdrawals {
		yooPayout, err := s.yookassa.GetPayout(withdrawal.ExternalID)
		if err != nil {
			s.log.Error("failed to check payout status", "error", err, "withdrawal_id", withdrawal.ID)
			continue
		}

		if yooPayout.Status == "succeeded" {
			s.log.Info("payout succeeded, updating status", "withdrawal_id", withdrawal.ID)
			if err := s.repo.UpdateWithdrawalStatus(ctx, withdrawal.ExternalID, core.WithdrawalSucceeded); err != nil {
				s.log.Error("failed to update withdrawal status", "error", err)
			}
		} else if yooPayout.Status == "failed" {
			s.log.Warn("payout failed, refunding", "withdrawal_id", withdrawal.ID)
			if err := s.repo.UpdateWithdrawalStatus(ctx, withdrawal.ExternalID, core.WithdrawalFailed); err != nil {
				s.log.Error("failed to update withdrawal status", "error", err)
				continue
			}

			if err := s.txClient.Deposit(ctx, withdrawal.EntityType, withdrawal.EntityID, withdrawal.Amount); err != nil {
				s.log.Error("failed to refund after failed payout", "error", err, "withdrawal_id", withdrawal.ID)
			}
		}
	}

	return nil
}
