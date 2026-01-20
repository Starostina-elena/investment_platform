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
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, p *core.Payment) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO payments (id, external_id, amount, entity_id, entity_type, status, created_at, updated_at)
		VALUES (:id, :external_id, :amount, :entity_id, :entity_type, :status, :created_at, :updated_at)
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

func (r *Repo) GetByID(ctx context.Context, id string) (*core.Payment, error) {
	var p core.Payment
	err := r.db.GetContext(ctx, &p, "SELECT * FROM payments WHERE id = $1", id)
	return &p, err
}

func (r *Repo) GetPendingPayments(ctx context.Context) ([]core.Payment, error) {
	var payments []core.Payment
	err := r.db.SelectContext(ctx, &payments, "SELECT * FROM payments WHERE status = $1", core.StatusPending)
	return payments, err
}

func (r *Repo) CreateWithdrawal(ctx context.Context, w *core.Withdrawal) error {
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO withdrawals (id, external_id, entity_id, entity_type, amount, status, created_at, updated_at)
		VALUES (:id, :external_id, :entity_id, :entity_type, :amount, :status, :created_at, :updated_at)
	`, w)
	return err
}

func (r *Repo) UpdateWithdrawalStatus(ctx context.Context, externalID string, status core.WithdrawalStatus) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE withdrawals SET status = $1, updated_at = NOW() WHERE external_id = $2
	`, status, externalID)
	return err
}

func (r *Repo) GetWithdrawalByExternalID(ctx context.Context, externalID string) (*core.Withdrawal, error) {
	var w core.Withdrawal
	err := r.db.GetContext(ctx, &w, "SELECT * FROM withdrawals WHERE external_id = $1", externalID)
	return &w, err
}

func (r *Repo) GetWithdrawalByID(ctx context.Context, id string) (*core.Withdrawal, error) {
	var w core.Withdrawal
	err := r.db.GetContext(ctx, &w, "SELECT * FROM withdrawals WHERE id = $1", id)
	return &w, err
}

func (r *Repo) GetPendingWithdrawals(ctx context.Context) ([]core.Withdrawal, error) {
	var withdrawals []core.Withdrawal
	err := r.db.SelectContext(ctx, &withdrawals, "SELECT * FROM withdrawals WHERE status = $1", core.WithdrawalPending)
	return withdrawals, err
}

func (r *Repo) GetWithdrawalsByEntity(ctx context.Context, entityType string, entityID int) ([]core.Withdrawal, error) {
	var withdrawals []core.Withdrawal
	err := r.db.SelectContext(ctx, &withdrawals,
		"SELECT * FROM withdrawals WHERE entity_type = $1 AND entity_id = $2 ORDER BY created_at DESC",
		entityType, entityID)
	return withdrawals, err
}
