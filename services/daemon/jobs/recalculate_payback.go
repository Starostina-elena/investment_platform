package jobs

import (
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type RecalculatePaybackJob struct {
	db  *sqlx.DB
	log *slog.Logger
}

type ProjectInfo struct {
	ID      int     `db:"id"`
	Percent float64 `db:"percent"`
}

type InvestorTransaction struct {
	UserID     int       `db:"user_id"`
	Amount     float64   `db:"amount"`
	InvestedAt time.Time `db:"invested_at"`
}

func NewRecalculatePaybackJob(db *sqlx.DB, log *slog.Logger) *RecalculatePaybackJob {
	return &RecalculatePaybackJob{
		db:  db,
		log: log,
	}
}

func (j *RecalculatePaybackJob) Run() {
	j.log.Info("starting recalculate payback job")

	var projects []ProjectInfo
	query := `
		SELECT id, percent
		FROM projects
		WHERE monetization_type = 'time_percent' AND payback_started = false
	`
	err := j.db.Select(&projects, query)
	if err != nil {
		j.log.Error("failed to fetch time_percent projects", "error", err)
		return
	}

	j.log.Info("found time_percent projects", "count", len(projects))

	for _, project := range projects {
		if err := j.recalculateForProject(project); err != nil {
			j.log.Error("failed to recalculate payback", "project_id", project.ID, "error", err)
		}
	}

	j.log.Info("recalculate payback job completed")
}

func (j *RecalculatePaybackJob) recalculateForProject(project ProjectInfo) error {
	investorPaybacks, err := j.calculateInvestorPaybacks(project.ID, project.Percent)
	if err != nil {
		return err
	}

	totalPayback := 0.0
	for _, payback := range investorPaybacks {
		totalPayback += payback
	}

	_, err = j.db.Exec(`
		UPDATE projects SET money_required_to_payback = $1 WHERE id = $2
	`, totalPayback, project.ID)
	if err != nil {
		j.log.Error("failed to update money_required_to_payback", "project_id", project.ID, "error", err)
		return err
	}

	j.log.Info("updated money_required_to_payback", "project_id", project.ID, "amount", totalPayback)
	return nil
}

func (j *RecalculatePaybackJob) calculateInvestorPaybacks(projectID int, percent float64) (map[int]float64, error) {
	var creatorID int
	err := j.db.Get(&creatorID, `SELECT creator_id FROM projects WHERE id = $1`, projectID)
	if err != nil {
		return nil, err
	}

	var investorTxs []InvestorTransaction
	query := `
		SELECT from_id as user_id, amount, time_at as invested_at
		FROM transactions
		WHERE reciever_id = $1 AND type = 'user_to_project' AND from_id != $2
		ORDER BY time_at ASC
	`
	err = j.db.Select(&investorTxs, query, projectID, creatorID)
	if err != nil {
		j.log.Error("failed to fetch investor transactions", "project_id", projectID, "error", err)
		return nil, err
	}

	investorMap := make(map[int][]InvestorTransaction)
	for _, tx := range investorTxs {
		investorMap[tx.UserID] = append(investorMap[tx.UserID], tx)
	}

	var usersWithPayback []int
	query = `
		SELECT DISTINCT reciever_id as user_id
		FROM transactions
		WHERE from_id = $1 AND type = 'project_to_user'
	`
	err = j.db.Select(&usersWithPayback, query, projectID)
	if err != nil {
		j.log.Error("failed to fetch users with payback", "project_id", projectID, "error", err)
		return nil, err
	}

	excludeUsers := make(map[int]bool)
	for _, userID := range usersWithPayback {
		excludeUsers[userID] = true
	}

	result := make(map[int]float64)
	now := time.Now()
	for userID, transactions := range investorMap {
		if excludeUsers[userID] {
			continue
		}

		paybackAmount := 0.0
		for _, tx := range transactions {
			days := int(now.Sub(tx.InvestedAt).Hours() / 24)
			paybackAmount += tx.Amount * (percent / 100) * float64(days)
		}
		result[userID] = paybackAmount
	}

	return result, nil
}
