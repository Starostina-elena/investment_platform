package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/organisation/service"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

func NewRepo(db *sqlx.DB, log slog.Logger) service.Repo {
	return &Repo{db: db, log: log}
}

func (r *Repo) Create(ctx context.Context, o *service.Org) (int, error) {
	var id int
	row := r.db.QueryRowxContext(ctx, `INSERT INTO organizations (name, owner, email) VALUES ($1,$2,$3) RETURNING id`, o.Name, 1, "")
	if err := row.Scan(&id); err != nil {
		r.log.Error("failed to insert organisation", "error", err)
		return 0, err
	}
	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*service.Org, error) {
	o := &service.Org{}
	if err := r.db.GetContext(ctx, o, `SELECT id, name FROM organizations WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		r.log.Error("failed to get organisation", "id", id, "error", err)
		return nil, err
	}
	return o, nil
}
