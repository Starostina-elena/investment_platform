package core

const (
	NotifTypeDividends = "dividends"
	NotifTypeRefund    = "refund"
)

type EmailRequest struct {
	Email       string  `json:"email"`
	Type        string  `json:"type"`
	ProjectName string  `json:"project_name"`
	Amount      float64 `json:"amount"`
}
