package repo

import (
	"context"
	"time"

	"github.com/Starostina-elena/investment_platform/services/payment/core"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	// Создаем таблицу при старте (для упрощения, в проде нужны миграции)
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS payments (
			id UUID PRIMARY KEY,
			external_id TEXT NOT NULL,
			amount NUMERIC(10, 2) NOT NULL,
			user_id INT NOT NULL,
			entity_type TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`)
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, p *core.Payment) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO payments (id, external_id, amount, user_id, entity_type, status, created_at, updated_at)
		VALUES (:id, :external_id, :amount, :user_id, :entity_type, :status, :created_at, :updated_at)
	`, p)
	return err
}

func (r *Repo) UpdateStatus(ctx context.Context, externalID string, status core.PaymentStatus) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE payments SET status = $1, updated_at = NOW() WHERE external_id = $2
	`, status, externalID)
	return err
}

func (r *Repo) GetByExternalID(ctx context.Context, externalID string) (*core.Payment, error) {
	var p core.Payment
	err := r.db.GetContext(ctx, &p, "SELECT * FROM payments WHERE external_id = $1", externalID)
	return &p, err
}
