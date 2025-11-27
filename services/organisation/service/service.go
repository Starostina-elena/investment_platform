package service

import (
	"context"
	"log/slog"
)

type Org struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Repo interface {
	Create(ctx context.Context, o *Org) (int, error)
	Get(ctx context.Context, id int) (*Org, error)
}
type Service interface {
	Create(ctx context.Context, name string) (*Org, error)
	Get(ctx context.Context, id int) (*Org, error)
}

type service struct {
	repo Repo
	log  slog.Logger
}

func NewService(r Repo, log slog.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, name string) (*Org, error) {
	o := &Org{Name: name}
	id, err := s.repo.Create(ctx, o)
	if err != nil {
		s.log.Error("failed to create organisation", "error", err)
		return nil, err
	}
	o.ID = id
	return o, nil
}

func (s *service) Get(ctx context.Context, id int) (*Org, error) {
	return s.repo.Get(ctx, id)
}
