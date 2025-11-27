package repo

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/transactions/service"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

func NewRepo(db *sqlx.DB, log slog.Logger) service.Repo { return &Repo{db: db, log: log} }

func (r *Repo) Create(ctx context.Context, t *service.Transaction) (int, error) {
	txType := "user_to_project"
	var id int
	row := r.db.QueryRowxContext(ctx, `INSERT INTO transactions (from_id, reciever_id, type, amount) VALUES ($1,$2,$3,$4) RETURNING id`, t.From, t.To, txType, t.Amount)
	if err := row.Scan(&id); err != nil {
		r.log.Error("failed to insert tx", "error", err)
		return 0, err
	}
	return id, nil
}
