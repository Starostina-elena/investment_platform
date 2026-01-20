package repo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Starostina-elena/investment_platform/services/transactions/clients"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Transaction struct {
	ID        int                `json:"id"`
	FromType  clients.EntityType `json:"from_type"`
	FromID    int                `json:"from_id"`
	ToType    clients.EntityType `json:"to_type"`
	ToID      int                `json:"to_id"`
	Amount    float64            `json:"amount"`
	CreatedAt time.Time          `json:"created_at"`
}

type Investor struct {
	UserID    int    `db:"user_id"`
	UserEmail string `db:"user_email"`
}

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

func NewRepo(db *sqlx.DB, log slog.Logger) *Repo {
	return &Repo{db: db, log: log}
}

func (r *Repo) Create(ctx context.Context, t *Transaction) (int, error) {
	txType := fmt.Sprintf("%s_to_%s", t.FromType, t.ToType)

	if t.FromID == 0 {
		txType = fmt.Sprintf("%s_deposit", t.ToType)
	}

	var id int
	row := r.db.QueryRowxContext(ctx,
		`INSERT INTO transactions (from_id, reciever_id, type, amount, time_at) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		t.FromID, t.ToID, txType, t.Amount, t.CreatedAt)

	if err := row.Scan(&id); err != nil {
		r.log.Error("failed to insert tx", "error", err)
		return 0, err
	}
	return id, nil
}

func (r *Repo) GetProjectInvestors(ctx context.Context, projectID int) ([]Investor, error) {
	var investors []Investor
	query := `
		SELECT DISTINCT ON (t.from_id) t.from_id as user_id, u.email as user_email
		FROM transactions t
		JOIN users u ON t.from_id = u.id
		WHERE t.reciever_id = $1 AND t.type = 'user_to_project'
	`
	err := r.db.SelectContext(ctx, &investors, query, projectID)
	if err != nil {
		r.log.Error("failed to get project investors", "error", err, "project_id", projectID)
		return nil, err
	}
	return investors, nil
}
