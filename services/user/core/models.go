package core

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Patronymic   *string   `json:"patronymic,omitempty"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	AvatarPath   *string   `json:"-" db:"avatar_path"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	IsBanned     bool      `json:"is_banned" db:"is_banned"`
}

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	Revoked   bool      `db:"revoked"`
}

type UserProjectInvestment struct {
	ProjectID        int       `json:"project_id" db:"project_id"`
	ProjectName      string    `json:"project_name" db:"project_name"`
	QuickPeek        string    `json:"quick_peek" db:"quick_peek"`
	MonetizationType string    `json:"monetization_type" db:"monetization_type"`
	TotalInvested    float64   `json:"total_invested" db:"total_invested"`
	TotalReceived    float64   `json:"total_received" db:"total_received"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	IsCompleted      bool      `json:"is_completed" db:"is_completed"`
	IsBanned         bool      `json:"is_banned" db:"is_banned"`
}
