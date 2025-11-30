package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Starostina-elena/investment_platform/services/user/auth"
	"github.com/Starostina-elena/investment_platform/services/user/core"
	"github.com/Starostina-elena/investment_platform/services/user/repo"
	"github.com/lib/pq"
)

type Service interface {
	Create(ctx context.Context, user core.User) (*core.User, error)
	Update(ctx context.Context, user core.User) (*core.User, error)
	Get(ctx context.Context, id int) (*core.User, error)
	GetByEmail(ctx context.Context, email string) (*core.User, error)
	GenerateRefreshToken(ctx context.Context, userID int, ttl time.Duration) (string, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	AuthenticateByRefresh(ctx context.Context, rawToken string) (*core.User, error)
}

type service struct {
	repo repo.RepoInterface
	log  slog.Logger
}

func NewService(r repo.RepoInterface, log slog.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, user core.User) (*core.User, error) {
	hashed, err := core.HashPassword(user.Password)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, err
	}
	user.PasswordHash = hashed
	user.Password = ""
	user.CreatedAt = time.Now()
	user.IsAdmin = false
	user.IsBanned = false
	user.Balance = 0.0
	id, err := s.repo.Create(ctx, &user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			s.log.Info("unique constraint violation", "constraint", pqErr.Constraint, "detail", pqErr.Detail)
			switch pqErr.Constraint {
			case "users_nickname_key":
				return nil, core.ErrNicknameExists
			case "users_email_key":
				return nil, core.ErrEmailExists
			}
		}
		s.log.Error("failed to create user", "error", err)
		return nil, err
	}
	user.ID = id
	return &user, nil
}

func (s *service) Update(ctx context.Context, user core.User) (*core.User, error) {
	hashed, err := core.HashPassword(user.Password)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, err
	}
	user.PasswordHash = hashed
	user.Password = ""
	user_updated, err := s.repo.Update(ctx, user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			s.log.Info("unique constraint violation", "constraint", pqErr.Constraint, "detail", pqErr.Detail)
			switch pqErr.Constraint {
			case "users_nickname_key":
				return nil, core.ErrNicknameExists
			case "users_email_key":
				return nil, core.ErrEmailExists
			}
		}
		s.log.Error("failed to update user", "error", err)
		return nil, err
	}
	return user_updated, nil
}

func (s *service) Get(ctx context.Context, id int) (*core.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user by email", "email", email, "error", err)
		return nil, err
	}
	return u, nil
}

func (s *service) GenerateRefreshToken(ctx context.Context, userID int, ttl time.Duration) (string, error) {
	raw, err := auth.GenerateRandomToken(32)
	if err != nil {
		s.log.Error("failed to generate refresh token", "error", err)
		return "", err
	}
	hash := auth.HashToken(raw)
	expiresAt := time.Now().Add(ttl)
	if _, err := s.repo.CreateRefreshToken(ctx, userID, hash, expiresAt); err != nil {
		s.log.Error("failed to store refresh token", "error", err)
		return "", err
	}
	return raw, nil
}

func (s *service) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	rt, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return err
	}
	return s.repo.RevokeRefreshToken(ctx, rt.ID)
}

func (s *service) AuthenticateByRefresh(ctx context.Context, rawToken string) (*core.User, error) {
	hash := auth.HashToken(rawToken)
	rt, err := s.repo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		s.log.Info("refresh: token not found", "error", err)
		return nil, err
	}
	if rt.Revoked {
		return nil, errors.New("token revoked")
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	u, err := s.repo.Get(ctx, rt.UserID)
	if err != nil {
		s.log.Error("failed to load user by refresh token", "user", rt.UserID, "error", err)
		return nil, err
	}
	return u, nil
}
