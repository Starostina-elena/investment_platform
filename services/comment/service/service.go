package service

import (
	"context"
	"log/slog"

	"github.com/Starostina-elena/investment_platform/services/comment/core"
	"github.com/Starostina-elena/investment_platform/services/comment/repo"
)

type Service interface {
	Create(ctx context.Context, projectID int, body string, userID int) (*core.Comment, error)
	Get(ctx context.Context, id int) (*core.Comment, error)
	Update(ctx context.Context, id int, body string) (*core.Comment, error)
	Delete(ctx context.Context, id int) error
	GetByProject(ctx context.Context, projectID int, limit, offset int) ([]core.Comment, error)
}

type service struct {
	repo repo.RepoInterface
	log  slog.Logger
}

func NewService(r repo.RepoInterface, log slog.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, projectID int, body string, userID int) (*core.Comment, error) {
	c := &core.Comment{ProjectID: projectID, Body: body, UserID: userID}
	id, err := s.repo.Create(ctx, c)
	if err != nil {
		s.log.Error("failed to create comment", "error", err)
		return nil, err
	}
	c.ID = id
	return c, nil
}

func (s *service) Get(ctx context.Context, id int) (*core.Comment, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id int, body string) (*core.Comment, error) {
	return s.repo.Update(ctx, id, body)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetByProject(ctx context.Context, projectID int, limit, offset int) ([]core.Comment, error) {
	return s.repo.GetByProject(ctx, projectID, limit, offset)
}
