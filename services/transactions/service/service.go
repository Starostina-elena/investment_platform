package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Starostina-elena/investment_platform/services/transactions/clients"
	"github.com/Starostina-elena/investment_platform/services/transactions/repo"
)

type Transaction = repo.Transaction

type Repo interface {
	Create(ctx context.Context, t *Transaction) (int, error)
	GetProjectInvestors(ctx context.Context, projectID int) ([]repo.Investor, error)
}

type Service interface {
	Transfer(ctx context.Context, fromType, toType clients.EntityType, fromID, toID int, amount float64) (*Transaction, error)
}

type service struct {
	repo               Repo
	balanceClient      *clients.BalanceClient
	projectClient      *clients.ProjectClient
	notificationClient *clients.NotificationClient
	log                slog.Logger
}

func NewService(repo Repo, bc *clients.BalanceClient, pc *clients.ProjectClient, nc *clients.NotificationClient, log slog.Logger) Service {
	return &service{
		repo:               repo,
		balanceClient:      bc,
		projectClient:      pc,
		notificationClient: nc,
		log:                log,
	}
}

func (s *service) Transfer(ctx context.Context, fromType, toType clients.EntityType, fromID, toID int, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	s.log.Info("starting transfer", "from", fromType, "from_id", fromID, "to", toType, "to_id", toID, "amount", amount)

	err := s.balanceClient.ChangeBalance(ctx, fromType, fromID, -amount)
	if err != nil {
		s.log.Error("failed to deduct funds", "error", err)
		return nil, fmt.Errorf("transaction failed: %v", err)
	}

	err = s.balanceClient.ChangeBalance(ctx, toType, toID, amount)
	if err != nil {
		s.log.Error("failed to add funds, starting rollback", "error", err)

		rbErr := s.balanceClient.ChangeBalance(ctx, fromType, fromID, amount)
		if rbErr != nil {
			s.log.Error("CRITICAL: ROLLBACK FAILED", "from_type", fromType, "from_id", fromID, "amount", amount, "error", rbErr)
			return nil, fmt.Errorf("system error: money stuck, contact support")
		}

		return nil, fmt.Errorf("transaction failed at destination: %v", err)
	}

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
		s.log.Error("transaction successful but failed to save record", "error", err)
	}
	t.ID = id

	if toType == clients.TypeProject {
		s.log.Info("processing project payment", "project_id", toID, "amount", amount)
		s.handleProjectPayment(ctx, toID, amount)
	}

	return t, nil
}

func (s *service) handleProjectPayment(ctx context.Context, projectID int, amount float64) {
	project, err := s.projectClient.GetProject(ctx, projectID)
	if err != nil {
		s.log.Error("failed to get project data", "error", err, "project_id", projectID)
		return
	}

	s.log.Info("got project data", "project_id", projectID, "monetization_type", project.MonetizationType)

	var paybackDelta float64
	switch project.MonetizationType {
	case "fixed_percent":
		paybackDelta = amount * (1 + project.Percent/100)
		s.log.Info("fixed_percent payback calculation", "amount", amount, "percent", project.Percent, "payback_delta", paybackDelta)

	case "time_percent":
		paybackDelta = amount
		s.log.Info("time_percent payback calculation", "amount", amount, "payback_delta", paybackDelta)

	default:
		paybackDelta = amount
		s.log.Info("default payback calculation", "monetization_type", project.MonetizationType, "payback_delta", paybackDelta)
	}

	newMoneyRequired := project.MoneyRequiredToPayback + paybackDelta
	if err := s.projectClient.UpdateMoneyRequiredToPayback(ctx, projectID, newMoneyRequired); err != nil {
		s.log.Error("failed to update money required to payback", "error", err, "project_id", projectID)
	} else {
		s.log.Info("updated money required to payback", "project_id", projectID, "new_amount", newMoneyRequired)
	}

	updatedProject, err := s.projectClient.GetProject(ctx, projectID)
	if err != nil {
		s.log.Error("failed to get updated project data for goal check", "error", err, "project_id", projectID)
		return
	}

	if updatedProject.CurrentMoney >= updatedProject.WantedMoney {
		s.log.Info("project goal reached!", "project_id", projectID, "current_money", updatedProject.CurrentMoney, "wanted_money", updatedProject.WantedMoney)
		s.sendGoalReachedNotifications(ctx, projectID, updatedProject)
	}
}

func (s *service) sendGoalReachedNotifications(ctx context.Context, projectID int, project *clients.ProjectData) {
	investors, err := s.repo.GetProjectInvestors(ctx, projectID)
	if err != nil {
		s.log.Error("failed to get project investors", "error", err, "project_id", projectID)
		return
	}

	s.log.Info("sending goal reached notifications", "project_id", projectID, "investor_count", len(investors))

	for _, investor := range investors {
		if err := s.notificationClient.SendEmail(investor.UserEmail, "project_goal_reached", project.Name, 0); err != nil {
			s.log.Error("failed to send goal reached email", "error", err, "user_id", investor.UserID, "email", investor.UserEmail)
		} else {
			s.log.Info("sent goal reached email", "user_id", investor.UserID, "email", investor.UserEmail)
		}
	}

	s.log.Info("goal reached notifications sent for project", "project_id", projectID, "project_name", project.Name)
}
