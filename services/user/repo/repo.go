package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/user/service"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

func NewRepo(db *sqlx.DB, log slog.Logger) service.Repo {
	return &Repo{db: db, log: log}
}

func (r *Repo) Create(ctx context.Context, u *service.User) (int, error) {
	surname := u.Name
	passwordHash := ""

	var id int
	row := r.db.QueryRowxContext(ctx,
		`INSERT INTO users (name, surname, nickname, email, password_hash) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		u.Name, surname, u.Nickname, u.Email, passwordHash,
	)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*service.User, error) {
	u := &service.User{}
	if err := r.db.GetContext(ctx, u, `SELECT id, name, nickname, email FROM users WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return u, nil
}
