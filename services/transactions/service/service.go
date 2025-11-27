package service

import (
	"context"
	"log/slog"
)

type Transaction struct {
	ID     int     `json:"id"`
	From   int     `json:"from"`
	To     int     `json:"to"`
	Amount float64 `json:"amount"`
}
type Repo interface {
	Create(cTransaction context.Context, t *Transaction) (int, error)
}
type Service interface {
	Create(cTransaction context.Context, from, to int, amount float64) (*Transaction, error)
}
type service struct {
	repo Repo
	log  slog.Logger
}

func NewService(repo Repo, log slog.Logger) Service { return &service{repo: repo, log: log} }
func (s *service) Create(cTransaction context.Context, from, to int, amount float64) (*Transaction, error) {
	t := &Transaction{From: from, To: to, Amount: amount}
	id, err := s.repo.Create(cTransaction, t)
	if err != nil {
		s.log.Error("failed to create Transaction", "error", err)
		return nil, err
	}
	t.ID = id
	return t, nil
}
