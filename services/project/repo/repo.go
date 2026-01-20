package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/project/core"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

type RepoInterface interface {
	Create(ctx context.Context, p *core.Project) (int, error)
	Get(ctx context.Context, id int) (*core.Project, error)
	Update(ctx context.Context, p *core.Project) (*core.Project, error)
	GetList(ctx context.Context, limit, offset int, monetizationType string) ([]core.Project, error)
	GetByCreator(ctx context.Context, creatorID int) ([]core.Project, error)
	GetAllByCreator(ctx context.Context, creatorID int) ([]core.Project, error)
	UpdatePicturePath(ctx context.Context, projectID int, picturePath *string) error
	BanProject(ctx context.Context, projectID int, banned bool) error
	ChangeProjectPublicity(ctx context.Context, projectID int, isPublic bool) error
	MarkProjectCompleted(ctx context.Context, projectID int, completed bool) error
	StartPayback(ctx context.Context, projectID int) error
	GetProjectTransactions(ctx context.Context, projectID int) ([]core.Transaction, error)
	UpdateMoneyRequiredToPayback(ctx context.Context, projectID int, newAmount float64) error
	AddFunds(ctx context.Context, projectID int, amount float64) error
}

func NewRepo(db *sqlx.DB, log slog.Logger) RepoInterface {
	return &Repo{db: db, log: log}
}

func (r *Repo) Create(ctx context.Context, p *core.Project) (int, error) {
	var id int
	row := r.db.QueryRowxContext(ctx,
		`INSERT INTO projects (name, creator_id, quick_peek, content, wanted_money, duration_days, is_public, monetization_type, percent) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		p.Name, p.CreatorID, p.QuickPeek, p.Content, p.WantedMoney, p.DurationDays, p.IsPublic, p.MonetizationType, p.Percent,
	)
	if err := row.Scan(&id); err != nil {
		r.log.Error("failed to insert project", "error", err)
		return 0, err
	}
	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*core.Project, error) {
	p := &core.Project{}
	if err := r.db.GetContext(ctx, p, `
		SELECT id, name, creator_id, quick_peek, quick_peek_picture_path, content, 
		       is_public, is_completed, current_money, wanted_money, duration_days, 
		       created_at, is_banned, monetization_type, percent, payback_started,
		       payback_started_date, money_required_to_payback
		FROM projects WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrProjectNotFound
		}
		r.log.Error("failed to get project", "id", id, "error", err)
		return nil, err
	}
	return p, nil
}

func (r *Repo) Update(ctx context.Context, p *core.Project) (*core.Project, error) {
	_, err := r.db.ExecContext(ctx,
		`UPDATE projects SET name=$1, quick_peek=$2, content=$3, is_public=$4, 
		wanted_money=$5, duration_days=$6 WHERE id=$7`,
		p.Name, p.QuickPeek, p.Content, p.IsPublic, p.WantedMoney, p.DurationDays, p.ID,
	)
	if err != nil {
		r.log.Error("failed to update project", "id", p.ID, "error", err)
		return nil, err
	}
	return r.Get(ctx, p.ID)
}

func (r *Repo) GetList(ctx context.Context, limit, offset int, monetizationType string) ([]core.Project, error) {
	projects := []core.Project{}

	query := `SELECT id, name, creator_id, quick_peek, quick_peek_picture_path, content, 
	       is_public, is_completed, current_money, wanted_money, duration_days,
	       payback_started_date, money_required_to_payback, 
	       created_at, is_banned, monetization_type, percent, payback_started
	FROM projects
	WHERE is_public = true AND is_banned = false AND is_completed = false`

	var args []interface{}
	argIndex := 1

	if monetizationType != "" {
		query += " AND monetization_type = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, monetizationType)
		argIndex++
	}

	query += " ORDER BY created_at DESC, id ASC LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, limit, offset)

	if err := r.db.SelectContext(ctx, &projects, query, args...); err != nil {
		r.log.Error("failed to get projects list", "error", err)
		return nil, err
	}

	return projects, nil
}

func (r *Repo) GetByCreator(ctx context.Context, creatorID int) ([]core.Project, error) {
	projects := []core.Project{}
	if err := r.db.SelectContext(ctx, &projects, `
		SELECT id, name, creator_id, quick_peek, quick_peek_picture_path, content, 
		       is_public, is_completed, current_money, wanted_money, duration_days,
		       payback_started_date, money_required_to_paybacks, 
		       created_at, is_banned, monetization_type, percent, payback_started
		FROM projects WHERE creator_id = $1 AND is_banned = false AND is_public = true ORDER BY created_at DESC, id ASC`, creatorID); err != nil {
		r.log.Error("failed to get projects by creator", "creator_id", creatorID, "error", err)
		return nil, err
	}
	return projects, nil
}

func (r *Repo) GetAllByCreator(ctx context.Context, creatorID int) ([]core.Project, error) {
	projects := []core.Project{}
	if err := r.db.SelectContext(ctx, &projects, `
		SELECT id, name, creator_id, quick_peek, quick_peek_picture_path, content,
		       payback_started_date, money_required_to_payback, 
		       is_public, is_completed, current_money, wanted_money, duration_days, 
		       created_at, is_banned, monetization_type, percent, payback_started
		FROM projects WHERE creator_id = $1 ORDER BY created_at DESC, id ASC`, creatorID); err != nil {
		r.log.Error("failed to get all projects by creator", "creator_id", creatorID, "error", err)
		return nil, err
	}
	return projects, nil
}

func (r *Repo) UpdatePicturePath(ctx context.Context, projectID int, picturePath *string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE projects SET quick_peek_picture_path = $1 WHERE id = $2`,
		picturePath, projectID,
	)
	if err != nil {
		r.log.Error("failed to update picture path", "project_id", projectID, "error", err)
		return err
	}
	return nil
}

func (r *Repo) BanProject(ctx context.Context, projectID int, banned bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE projects SET is_banned = $1 WHERE id = $2`,
		banned, projectID,
	)
	if err != nil {
		r.log.Error("failed to ban/unban project", "project_id", projectID, "banned", banned, "error", err)
		return err
	}
	return nil
}

func (r *Repo) ChangeProjectPublicity(ctx context.Context, projectID int, isPublic bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE projects SET is_public = $1 WHERE id = $2`,
		isPublic, projectID,
	)
	if err != nil {
		r.log.Error("failed to change project publicity", "project_id", projectID, "is_public", isPublic, "error", err)
		return err
	}
	return nil
}

func (r *Repo) MarkProjectCompleted(ctx context.Context, projectID int, completed bool) error {
	var paybackStarted bool
	err := r.db.GetContext(ctx, &paybackStarted, `SELECT payback_started FROM projects WHERE id = $1`, projectID)
	if err != nil {
		r.log.Error("failed to check payback_started status", "project_id", projectID, "error", err)
		return err
	}

	if paybackStarted {
		r.log.Warn("cannot change is_completed when payback_started is true", "project_id", projectID)
		return core.ErrPaybackStarted
	}

	_, err = r.db.ExecContext(ctx,
		`UPDATE projects SET is_completed = $1 WHERE id = $2`,
		completed, projectID,
	)
	if err != nil {
		r.log.Error("failed to mark project as completed/incomplete", "project_id", projectID, "completed", completed, "error", err)
		return err
	}
	return nil
}

func (r *Repo) StartPayback(ctx context.Context, projectID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE projects SET payback_started = true,
		payback_started_date = CURRENT_TIMESTAMP, is_completed = true WHERE id = $1`,
		projectID,
	)
	if err != nil {
		r.log.Error("failed to start payback", "project_id", projectID, "error", err)
		return err
	}
	r.log.Info("payback started", "project_id", projectID)
	return nil
}

func (r *Repo) GetProjectTransactions(ctx context.Context, projectID int) ([]core.Transaction, error) {
	var transactions []core.Transaction
	err := r.db.SelectContext(ctx, &transactions,
		`SELECT id, from_id, reciever_id, type, amount, time_at 
		FROM transactions 
		WHERE (from_id = $1 AND type = 'project_to_user') 
		   OR (reciever_id = $1 AND type = 'user_to_project')
		ORDER BY time_at ASC`,
		projectID,
	)
	if err != nil {
		r.log.Error("failed to get project transactions", "project_id", projectID, "error", err)
		return nil, err
	}
	return transactions, nil
}

func (r *Repo) UpdateMoneyRequiredToPayback(ctx context.Context, projectID int, newAmount float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE projects SET money_required_to_payback = $1 WHERE id = $2`,
		newAmount, projectID,
	)
	if err != nil {
		r.log.Error("failed to update money_required_to_payback", "project_id", projectID, "error", err)
		return err
	}
	return nil
}

func (r *Repo) AddFunds(ctx context.Context, projectID int, amount float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE projects SET current_money = current_money + $1 WHERE id = $2`,
		amount, projectID)
	if err != nil {
		r.log.Error("failed to add funds", "project_id", projectID, "amount", amount, "error", err)
		return err
	}
	return nil
}
