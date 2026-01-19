package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/comment/cache"
	"github.com/Starostina-elena/investment_platform/services/comment/core"
)

type RepoInterface interface {
	Create(ctx context.Context, c *core.Comment) (int, error)
	Get(ctx context.Context, id int) (*core.Comment, error)
	Update(ctx context.Context, id int, body string) (*core.Comment, error)
	Delete(ctx context.Context, id int) error
	GetByProject(ctx context.Context, projectID int, limit, offset int) ([]core.Comment, error)
}

type Repo struct {
	db    *sqlx.DB
	cache *cache.Cache
	log   slog.Logger
}

func NewRepo(db *sqlx.DB, c *cache.Cache, log slog.Logger) RepoInterface {
	return &Repo{db: db, cache: c, log: log}
}

func (r *Repo) Create(ctx context.Context, c *core.Comment) (int, error) {
	var id int
	row := r.db.QueryRowxContext(ctx, `INSERT INTO comments (body, user_id, project_id, created_at) VALUES ($1,$2,$3,NOW()) RETURNING id`, c.Body, c.UserID, c.ProjectID)
	if err := row.Scan(&id); err != nil {
		r.log.Error("failed to insert comment", "error", err)
		return 0, err
	}
	c.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.000000Z")

	var username string
	err := r.db.GetContext(ctx, &username, `SELECT nickname FROM users WHERE id = $1`, c.UserID)
	if err != nil {
		r.log.Error("failed to get username", "user_id", c.UserID, "error", err)
		username = ""
	}

	c.ID = id
	c.Username = username
	_ = r.cache.SetComment(ctx, c)
	_ = r.cache.InvalidateProjectComments(ctx, c.ProjectID)
	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*core.Comment, error) {
	if cached, err := r.cache.GetComment(ctx, id); err == nil && cached != nil {
		return cached, nil
	}

	c := &core.Comment{}
	if err := r.db.GetContext(ctx, c,
		`SELECT c.id, c.project_id, c.user_id, u.nickname as username, c.body, c.created_at 
		 FROM comments c 
		 JOIN users u ON c.user_id = u.id 
		 WHERE c.id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrCommentNotFound
		}
		r.log.Error("failed to get comment", "id", id, "error", err)
		return nil, err
	}
	_ = r.cache.SetComment(ctx, c)
	return c, nil
}

func (r *Repo) Update(ctx context.Context, id int, body string) (*core.Comment, error) {
	_, err := r.db.ExecContext(ctx, `UPDATE comments SET body = $1 WHERE id = $2`, body, id)
	if err != nil {
		r.log.Error("failed to update comment", "id", id, "error", err)
		return nil, err
	}
	_ = r.cache.DeleteComment(ctx, id)

	return r.Get(ctx, id)
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	c, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		r.log.Error("failed to delete comment", "id", id, "error", err)
		return err
	}
	_ = r.cache.DeleteComment(ctx, id)
	_ = r.cache.InvalidateProjectComments(ctx, c.ProjectID)
	return nil
}

func (r *Repo) GetByProject(ctx context.Context, projectID int, limit, offset int) ([]core.Comment, error) {
	if cached, err := r.cache.GetProjectComments(ctx, projectID, limit, offset); err == nil && cached != nil {
		return cached, nil
	}

	var comments []core.Comment
	err := r.db.SelectContext(ctx, &comments,
		`SELECT c.id, c.project_id, c.user_id, u.nickname as username, c.body, c.created_at 
		 FROM comments c 
		 JOIN users u ON c.user_id = u.id 
		 WHERE c.project_id = $1 
		 ORDER BY c.created_at DESC 
		 LIMIT $2 OFFSET $3`,
		projectID, limit, offset)
	if err != nil {
		r.log.Error("failed to get comments by project", "project_id", projectID, "error", err)
		return nil, err
	}
	_ = r.cache.SetProjectComments(ctx, projectID, limit, offset, comments)
	return comments, nil
}
