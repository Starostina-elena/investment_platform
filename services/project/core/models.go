package core

import "time"

type Project struct {
	ID                     int        `json:"id" db:"id"`
	Name                   string     `json:"name" db:"name"`
	CreatorID              int        `json:"creator_id" db:"creator_id"`
	QuickPeek              string     `json:"quick_peek" db:"quick_peek"`
	QuickPeekPicturePath   *string    `json:"-" db:"quick_peek_picture_path"`
	Content                string     `json:"content" db:"content"`
	IsPublic               bool       `json:"is_public" db:"is_public"`
	IsCompleted            bool       `json:"is_completed" db:"is_completed"`
	CurrentMoney           float64    `json:"current_money" db:"current_money"`
	WantedMoney            float64    `json:"wanted_money" db:"wanted_money"`
	DurationDays           int        `json:"duration_days" db:"duration_days"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	IsBanned               bool       `json:"is_banned" db:"is_banned"`
	MonetizationType       string     `json:"monetization_type" db:"monetization_type"`
	Percent                float64    `json:"percent,omitempty" db:"percent"`
	PaybackStarted         bool       `json:"payback_started" db:"payback_started"`
	PaybackStartedDate     *time.Time `json:"payback_started_date,omitempty" db:"payback_started_date"`
	MoneyRequiredToPayback float64    `json:"money_required_to_payback" db:"money_required_to_payback"`
}

type Transaction struct {
	ID             int       `db:"id"`
	FromID         *int      `db:"from_id"`
	ReceiverID     *int      `db:"reciever_id"`
	Type           string    `db:"type"`
	Amount         float64   `db:"amount"`
	TimeAt         time.Time `db:"time_at"`
}

type InvestorPayback struct {
	UserID         int
	TotalInvested  float64
	TotalReceived  float64
	PaybackAmount  float64
	Investments    []Investment
}

type Investment struct {
	Amount    float64
	InvestedAt time.Time
}
