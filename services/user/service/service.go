package service

import (
	"context"
	"log/slog"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type Repo interface {
	Create(ctx context.Context, u *User) (int, error)
	Get(ctx context.Context, id int) (*User, error)
}

type Service interface {
	Create(ctx context.Context, name, nickname, email string) (*User, error)
	Get(ctx context.Context, id int) (*User, error)
}

type service struct {
	repo Repo
	log  slog.Logger
}

func NewService(r Repo, log slog.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, name, nickname, email string) (*User, error) {
	u := &User{Name: name, Nickname: nickname, Email: email}
	id, err := s.repo.Create(ctx, u)
	if err != nil {
		s.log.Error("failed to create user", "error", err)
		return nil, err
	}
	u.ID = id
	return u, nil
}

func (s *service) Get(ctx context.Context, id int) (*User, error) {
	return s.repo.Get(ctx, id)
}
