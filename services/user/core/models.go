package core

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Patronymic   string    `json:"patronymic,omitempty"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	AvatarPath   string    `json:"avatar_path,omitempty"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	IsAdmin      bool      `json:"is_admin"`
	IsBanned     bool      `json:"is_banned"`
}
