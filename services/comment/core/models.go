package core

type Comment struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	Username  string `json:"username" db:"username"`
	ProjectID int    `json:"project_id" db:"project_id"`
	Body      string `json:"body" db:"body"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
