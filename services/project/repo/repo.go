package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/project/service"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

func NewRepo(db *sqlx.DB, log slog.Logger) service.Repo { return &Repo{db: db, log: log} }

func (r *Repo) Create(ctx context.Context, p *service.Project) (int, error) {
	creatorID := 1
	quickPeek := ""
	content := ""
	wantedMoney := 1.0
	var id int
	row := r.db.QueryRowxContext(ctx, `INSERT INTO projects (name, creator_id, quick_peek, content, wanted_money) VALUES ($1,$2,$3,$4,$5) RETURNING id`, p.Name, creatorID, quickPeek, content, wantedMoney)
	if err := row.Scan(&id); err != nil {
		r.log.Error("failed to insert project", "error", err)
		return 0, err
	}
	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*service.Project, error) {
	p := &service.Project{}
	if err := r.db.GetContext(ctx, p, `SELECT id, name FROM projects WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		r.log.Error("failed to get project", "id", id, "error", err)
		return nil, err
	}
	return p, nil
}
