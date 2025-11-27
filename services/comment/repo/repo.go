package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/comment/service"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

func NewRepo(db *sqlx.DB, log slog.Logger) service.Repo {
	return &Repo{db: db, log: log}
}

func (r *Repo) Create(ctx context.Context, c *service.Comment) (int, error) {
	var id int
	row := r.db.QueryRowxContext(ctx, `INSERT INTO comments (body, user_id, project_id, created_at) VALUES ($1,$2,$3,NOW()) RETURNING id`, c.Body, 1, c.ProjectID)
	if err := row.Scan(&id); err != nil {
		r.log.Error("failed to insert comment", "error", err)
		return 0, err
	}
	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*service.Comment, error) {
	c := &service.Comment{}
	if err := r.db.GetContext(ctx, c, `SELECT id, project_id, body FROM comments WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		r.log.Error("failed to get comment", "id", id, "error", err)
		return nil, err
	}
	return c, nil
}
