package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

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
	Get(ctx context.Context, id int) (*core.User, error)
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
