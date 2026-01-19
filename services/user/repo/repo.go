package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/user/core"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

type RepoInterface interface {
	Create(ctx context.Context, u *core.User) (int, error)
	Update(ctx context.Context, user core.User) (*core.User, error)
	Get(ctx context.Context, id int) (*core.User, error)
	GetByEmail(ctx context.Context, email string) (*core.User, error)
	CreateRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (int, error)
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*core.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id int) error
	RevokeAllRefreshTokens(ctx context.Context, userID int) error
	SetAdmin(ctx context.Context, userID int, isAdmin bool) error
	BanUser(ctx context.Context, userID int, isBanned bool) error
	UpdateAvatarPath(ctx context.Context, userID int, avatarPath *string) error
	UpdatePassword(ctx context.Context, userID int, passwordHash string) error
}

func NewRepo(db *sqlx.DB, log slog.Logger) RepoInterface {
	return &Repo{db: db, log: log}
}

func (r *Repo) Create(ctx context.Context, u *core.User) (int, error) {

	var id int
	row := r.db.QueryRowxContext(ctx,
		`INSERT INTO users (name, surname, patronymic, nickname, email, password_hash, balance, created_at, is_admin, is_banned) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		u.Name, u.Surname, u.Patronymic, u.Nickname, u.Email, u.PasswordHash, u.Balance, u.CreatedAt, u.IsAdmin, u.IsBanned,
	)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) Update(ctx context.Context, user core.User) (*core.User, error) {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET name=$1, surname=$2, patronymic=$3, nickname=$4, email=$5 WHERE id=$6`,
		user.Name, user.Surname, user.Patronymic, user.Nickname, user.Email, user.ID,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*core.User, error) {
	u := &core.User{}
	if err := r.db.GetContext(ctx, u, `SELECT id, name, surname, patronymic, nickname, email, avatar_path, password_hash, balance, created_at, is_admin, is_banned FROM users WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return u, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	u := &core.User{}
	if err := r.db.GetContext(ctx, u, `SELECT id, name, surname, patronymic, nickname, email, avatar_path, password_hash, balance, created_at, is_admin, is_banned FROM users WHERE email = $1`, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return u, nil
}

func (r *Repo) CreateRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (int, error) {
	var id int
	row := r.db.QueryRowxContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3) RETURNING id`,
		userID, tokenHash, expiresAt,
	)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*core.RefreshToken, error) {
	rt := &core.RefreshToken{}
	if err := r.db.GetContext(ctx, rt, `SELECT id, user_id, token_hash, expires_at, created_at, revoked FROM refresh_tokens WHERE token_hash = $1`, tokenHash); err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *Repo) RevokeRefreshToken(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked = true WHERE id = $1`, id)
	return err
}

func (r *Repo) RevokeAllRefreshTokens(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1`, userID)
	return err
}

func (r *Repo) SetAdmin(ctx context.Context, userID int, isAdmin bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET is_admin = $1 WHERE id = $2`, isAdmin, userID)
	return err
}

func (r *Repo) BanUser(ctx context.Context, userID int, isBanned bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET is_banned = $1 WHERE id = $2`, isBanned, userID)
	return err
}

func (r *Repo) UpdateAvatarPath(ctx context.Context, userID int, avatarPath *string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET avatar_path = $1 WHERE id = $2`, avatarPath, userID)
	return err
}

func (r *Repo) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET password_hash = $1 WHERE id = $2`, passwordHash, userID)
	return err
}
