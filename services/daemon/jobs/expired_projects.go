package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

type ExpiredProjectsJob struct {
	db              *sqlx.DB
	log             *slog.Logger
	notificationURL string
}

type Investment struct {
	UserID      int     `db:"user_id"`
	UserEmail   string  `db:"user_email"`
	Amount      float64 `db:"amount"`
	ProjectName string  `db:"project_name"`
}

type Project struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	OwnerEmail  string    `db:"owner_email"`
	CreatedAt   time.Time `db:"created_at"`
	DurationDays int      `db:"duration_days"`
}

func NewExpiredProjectsJob(db *sqlx.DB, log *slog.Logger) *ExpiredProjectsJob {
	notifURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notifURL == "" {
		notifURL = "http://notification:8083"
	}
	return &ExpiredProjectsJob{
		db:              db,
		log:             log,
		notificationURL: notifURL,
	}
}

func (j *ExpiredProjectsJob) Run() {
	j.log.Info("starting expired projects job")

	var expiredProjects []Project
	query := `
		SELECT p.id, p.name, u.email as owner_email, p.created_at, p.duration_days
		FROM projects p
		JOIN organizations o ON p.creator_id = o.id
		JOIN users u ON o.owner = u.id
		WHERE p.created_at + (p.duration_days || ' days')::interval < NOW() AND p.is_completed = false
	`
	err := j.db.Select(&expiredProjects, query)
	if err != nil {
		j.log.Error("failed to fetch expired projects", "error", err)
		return
	}

	j.log.Info("found expired projects", "count", len(expiredProjects))

	for _, project := range expiredProjects {
		if err := j.processExpiredProject(project); err != nil {
			j.log.Error("failed to process expired project", "project_id", project.ID, "error", err)
		}
	}

	j.log.Info("expired projects job completed")
}

func (j *ExpiredProjectsJob) processExpiredProject(project Project) error {
	tx, err := j.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var investments []Investment
	query := `
		SELECT DISTINCT ON (t.from_id) t.from_id as user_id, u.email as user_email, $1 as project_name, 0 as amount
		FROM transactions t
		JOIN users u ON t.from_id = u.id
		WHERE t.reciever_id = $2 AND t.type = 'user_to_project'
	`
	err = tx.Select(&investments, query, project.Name, project.ID)
	if err != nil {
		return fmt.Errorf("fetch investments: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE projects SET is_completed = true, is_public = false WHERE id = $1
	`, project.ID)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	for _, inv := range investments {
		j.sendEmail(inv.UserEmail, "project_closed", project.Name, 0)
	}

	j.sendEmailToOwner(project.OwnerEmail, project.Name)

	j.log.Info("processed expired project", "project_id", project.ID, "investors", len(investments))
	return nil
}

func (j *ExpiredProjectsJob) sendEmail(email, notifType, projectName string, amount float64) {
	payload := map[string]interface{}{
		"email":        email,
		"type":         notifType,
		"project_name": projectName,
		"amount":       amount,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(j.notificationURL+"/send", "application/json", bytes.NewReader(body))
	if err != nil {
		j.log.Error("failed to send email", "error", err, "to", email)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		j.log.Error("email service returned error", "status", resp.StatusCode, "to", email)
	}
}

func (j *ExpiredProjectsJob) sendEmailToOwner(email, projectName string) {
	j.sendEmail(email, "project_closed", projectName, 0)
}
