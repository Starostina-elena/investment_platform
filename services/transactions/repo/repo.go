package repo

import (
	"context"
	"log/slog"

	"github.com/Starostina-elena/investment_platform/services/transactions/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

func (r *Repo) Create(ctx context.Context, t *service.Transaction) (int, error) {
	// Используем type для обозначения направления, method можно логировать или добавить колонку в БД позже
	var id int
	row := r.db.QueryRowxContext(ctx,
		`INSERT INTO transactions (from_id, reciever_id, type, amount, time_at) 
         VALUES ($1,$2,$3,$4, NOW()) RETURNING id`,
		t.FromID, t.ToID, t.Type, t.Amount)
	if err := row.Scan(&id); err != nil {
		r.log.Error("failed to insert tx", "error", err)
		return 0, err
	}
	return id, nil
}

func NewRepo(db *sqlx.DB, log slog.Logger) service.Repo { return &Repo{db: db, log: log} }
