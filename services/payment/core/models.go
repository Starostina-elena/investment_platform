package core

import "time"

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusSucceeded PaymentStatus = "succeeded"
	StatusCanceled  PaymentStatus = "canceled"
)

type Payment struct {
	ID         string        `db:"id"`          // Внутренний UUID
	ExternalID string        `db:"external_id"` // ID в ЮKassa
	Amount     float64       `db:"amount"`
	UserID     int           `db:"user_id"`     // Кого пополняем
	EntityType string        `db:"entity_type"` // "user" или "org"
	Status     PaymentStatus `db:"status"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
}
