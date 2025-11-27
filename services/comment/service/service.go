package service

import (
	"context"
	"log/slog"
)

type Comment struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Body      string `json:"body"`
}
type Repo interface {
	Create(ctx context.Context, c *Comment) (int, error)
	Get(ctx context.Context, id int) (*Comment, error)
}
type Service interface {
	Create(ctx context.Context, projectID int, body string) (*Comment, error)
	Get(ctx context.Context, id int) (*Comment, error)
}
type service struct {
	repo Repo
	log  slog.Logger
}

func NewService(r Repo, log slog.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, projectID int, body string) (*Comment, error) {
	c := &Comment{ProjectID: projectID, Body: body}
	id, err := s.repo.Create(ctx, c)
	if err != nil {
		s.log.Error("failed to create comment", "error", err)
		return nil, err
	}
	c.ID = id
	return c, nil
}

func (s *service) Get(ctx context.Context, id int) (*Comment, error) {
	return s.repo.Get(ctx, id)
}
