package core

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Patronymic   *string   `json:"patronymic,omitempty"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	AvatarPath   *string   `json:"avatar_path,omitempty" db:"avatar_path"`
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
