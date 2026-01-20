package service

import (
	"context"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/Starostina-elena/investment_platform/services/project/clients"
	"github.com/Starostina-elena/investment_platform/services/project/core"
	"github.com/Starostina-elena/investment_platform/services/project/repo"
	"github.com/Starostina-elena/investment_platform/services/project/storage"
)

type Service interface {
	Create(ctx context.Context, req core.Project, creatorID int, userID int) (*core.Project, error)
	Get(ctx context.Context, id int) (*core.Project, error)
	Update(ctx context.Context, projectID int, p core.Project, userID int) (*core.Project, error)
	GetList(ctx context.Context, limit, offset int, monetizationType string) ([]core.Project, error)
	GetByCreator(ctx context.Context, creatorID int) ([]core.Project, error)
	GetAllByCreator(ctx context.Context, projectID int, userID int, isAdmin bool) ([]core.Project, error)
	UpdatePicturePath(ctx context.Context, projectID int, picturePath string) error
	BanProject(ctx context.Context, projectID int, banned bool) error
	ChangeProjectPublicity(ctx context.Context, projectID int, userID int, isPublic bool) error
	MarkProjectCompleted(ctx context.Context, projectID int, userID int, completed bool) error
	StartPayback(ctx context.Context, projectID int, userID int) error
	UploadPicture(ctx context.Context, projectID int, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	deletePicture(ctx context.Context, projectID int, picturePath string) error
	DeletePictureFromProject(ctx context.Context, projectID int, userID int) error
}

type service struct {
	repo      repo.RepoInterface
	orgClient *clients.OrgClient
	minio     *storage.MinioStorage
	log       slog.Logger
}

func NewService(r repo.RepoInterface, orgClient *clients.OrgClient, minioStorage *storage.MinioStorage, log slog.Logger) Service {
	return &service{repo: r, orgClient: orgClient, minio: minioStorage, log: log}
}

func (s *service) Create(ctx context.Context, p core.Project, creatorID int, userID int) (*core.Project, error) {
	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, creatorID, userID, "project_management")
	if err != nil {
		s.log.Error("failed to check organisation permission", "error", err)
		return nil, err
	}
	if !allowed {
		return nil, core.ErrNotAuthorized
	}

	p.IsPublic = true
	p.IsCompleted = false
	p.CurrentMoney = 0.0
	p.CreatedAt = time.Now()
	p.IsBanned = false

	id, err := s.repo.Create(ctx, &p)
	if err != nil {
		s.log.Error("failed to create project", "error", err)
		return nil, err
	}
	p.ID = id
	return &p, nil
}

func (s *service) Get(ctx context.Context, id int) (*core.Project, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, projectID int, p core.Project, userID int) (*core.Project, error) {
	existingProject, err := s.repo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}

	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, existingProject.CreatorID, userID, "project_management")
	if err != nil || !allowed {
		return nil, core.ErrNotAuthorized
	}

	existingProject.Name = p.Name
	existingProject.QuickPeek = p.QuickPeek
	existingProject.Content = p.Content
	existingProject.IsPublic = p.IsPublic
	existingProject.WantedMoney = p.WantedMoney
	existingProject.DurationDays = p.DurationDays

	updatedProject, err := s.repo.Update(ctx, existingProject)
	if err != nil {
		s.log.Error("failed to update project", "error", err)
		return nil, err
	}
	return updatedProject, nil
}

func (s *service) GetList(ctx context.Context, limit, offset int, monetizationType string) ([]core.Project, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetList(ctx, limit, offset, monetizationType)
}

func (s *service) GetByCreator(ctx context.Context, creatorID int) ([]core.Project, error) {
	return s.repo.GetByCreator(ctx, creatorID)
}

func (s *service) GetAllByCreator(ctx context.Context, orgID int, userID int, isAdmin bool) ([]core.Project, error) {
	if !isAdmin {
		allowed, err := s.orgClient.CheckUserOrgPermission(ctx, orgID, userID, "project_management")
		if err != nil || !allowed {
			return nil, core.ErrNotAuthorized
		}
	}
	return s.repo.GetAllByCreator(ctx, orgID)
}

func (s *service) UpdatePicturePath(ctx context.Context, projectID int, picturePath string) error {
	return s.repo.UpdatePicturePath(ctx, projectID, &picturePath)
}

func (s *service) BanProject(ctx context.Context, projectID int, banned bool) error {
	return s.repo.BanProject(ctx, projectID, banned)
}

func (s *service) ChangeProjectPublicity(ctx context.Context, projectID int, userID int, isPublic bool) error {
	existingProject, err := s.repo.Get(ctx, projectID)
	if err != nil {
		return err
	}
	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, existingProject.CreatorID, userID, "project_management")
	if err != nil || !allowed {
		return core.ErrNotAuthorized
	}
	return s.repo.ChangeProjectPublicity(ctx, projectID, isPublic)
}

func (s *service) MarkProjectCompleted(ctx context.Context, projectID int, userID int, completed bool) error {
	existingProject, err := s.repo.Get(ctx, projectID)
	if err != nil {
		return err
	}
	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, existingProject.CreatorID, userID, "project_management")
	if err != nil || !allowed {
		return core.ErrNotAuthorized
	}
	return s.repo.MarkProjectCompleted(ctx, projectID, completed)
}

func (s *service) StartPayback(ctx context.Context, projectID int, userID int) error {
	existingProject, err := s.repo.Get(ctx, projectID)
	if err != nil {
		return err
	}

	allowed, err := s.orgClient.CheckUserOrgPermission(ctx, existingProject.CreatorID, userID, "money_management")
	if err != nil {
		s.log.Error("failed to check money_management permission", "error", err)
		return err
	}
	if !allowed {
		return core.ErrNotAuthorized
	}

	if existingProject.PaybackStarted {
		s.log.Warn("payback already started", "project_id", projectID)
		return core.ErrPaybackStarted
	}

	if existingProject.MonetizationType == "charity" || existingProject.MonetizationType == "custom" {
		s.log.Warn("payback not supported for monetization type", "project_id", projectID, "type", existingProject.MonetizationType)
		return core.ErrPaybackNotSupported
	}

	err = s.repo.StartPayback(ctx, projectID)
	if err != nil {
		s.log.Error("failed to start payback", "error", err)
		return err
	}

	transactions, err := s.repo.GetProjectTransactions(ctx, projectID)
	if err != nil {
		s.log.Error("failed to get project transactions", "error", err)
		return err
	}

	investorPaybacks := s.calculatePaybacks(existingProject, transactions)
	newMoneyRequired := existingProject.MoneyRequiredToPayback
	defer func() {
		_ = s.repo.UpdateMoneyRequiredToPayback(ctx, projectID, newMoneyRequired)
	}()

	for _, payback := range investorPaybacks {
		if payback.TotalReceived > 0 {
			continue
		}
		amountToPay := payback.PaybackAmount

		if existingProject.CurrentMoney >= amountToPay {
			// TODO: вызов микросервиса транзакций для создания транзакции project_to_user
			// transactionService.CreateTransaction(ctx, projectID, payback.UserID, amountToPay, "project_to_user")
			existingProject.CurrentMoney -= amountToPay
			newMoneyRequired -= amountToPay
			if newMoneyRequired < 0 {
				newMoneyRequired = 0
			}
		} else {
			s.log.Info("insufficient funds for full payback", "project_id", projectID, "user_id", payback.UserID, "amount_needed", amountToPay, "current_money", existingProject.CurrentMoney)
			return core.ErrNotEnoughFunds
		}
	}

	return nil
}

func (s *service) calculatePaybacks(project *core.Project, transactions []core.Transaction) []core.InvestorPayback {
	investorMap := make(map[int]*core.InvestorPayback)

	for _, tx := range transactions {
		if tx.Type == "user_to_project" && tx.ReceiverID != nil && *tx.ReceiverID == project.ID && tx.FromID != nil {
			userID := *tx.FromID
			if userID == project.CreatorID {
				continue
			}
			if _, exists := investorMap[userID]; !exists {
				investorMap[userID] = &core.InvestorPayback{
					UserID:        userID,
					TotalInvested: 0,
					TotalReceived: 0,
					Investments:   []core.Investment{},
				}
			}
			investorMap[userID].TotalInvested += tx.Amount
			investorMap[userID].Investments = append(investorMap[userID].Investments, core.Investment{
				Amount:     tx.Amount,
				InvestedAt: tx.TimeAt,
			})
		} else if tx.Type == "project_to_user" && tx.FromID != nil && *tx.FromID == project.ID && tx.ReceiverID != nil {
			userID := *tx.ReceiverID
			if userID == project.CreatorID {
				continue
			}
			if _, exists := investorMap[userID]; !exists {
				investorMap[userID] = &core.InvestorPayback{
					UserID:        userID,
					TotalInvested: 0,
					TotalReceived: 0,
					Investments:   []core.Investment{},
				}
			}
			investorMap[userID].TotalReceived += tx.Amount
		}
	}

	result := make([]core.InvestorPayback, 0, len(investorMap))
	for _, investor := range investorMap {
		switch project.MonetizationType {
		case "fixed_percent":
			investor.PaybackAmount = investor.TotalInvested * (1 + project.Percent/100)
		case "time_percent":
			paybackAmount := 0.0
			now := time.Now()
			for _, inv := range investor.Investments {
				days := int(now.Sub(inv.InvestedAt).Hours() / 24)
				paybackAmount += inv.Amount * (project.Percent / 100) * float64(days)
			}
			investor.PaybackAmount = paybackAmount
		}
		result = append(result, *investor)
	}

	return result
}
