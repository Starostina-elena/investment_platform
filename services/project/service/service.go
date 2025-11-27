package service

import (
	"context"
	"log/slog"
)

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Repo interface {
	Create(ctx context.Context, p *Project) (int, error)
	Get(ctx context.Context, id int) (*Project, error)
}
type Service interface {
	Create(ctx context.Context, name string) (*Project, error)
	Get(ctx context.Context, id int) (*Project, error)
}
type service struct {
	repo Repo
	log  slog.Logger
}

func NewService(r Repo, log slog.Logger) Service { return &service{repo: r, log: log} }
func (s *service) Create(ctx context.Context, name string) (*Project, error) {
	p := &Project{Name: name}
	id, err := s.repo.Create(ctx, p)
	if err != nil {
		s.log.Error("failed to create project", "error", err)
		return nil, err
	}
	p.ID = id
	return p, nil
}
func (s *service) Get(ctx context.Context, id int) (*Project, error) { return s.repo.Get(ctx, id) }
