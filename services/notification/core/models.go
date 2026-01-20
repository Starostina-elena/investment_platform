package core

const (
	NotifTypeDividends     = "dividends"
	NotifTypeProjectClosed = "project_closed"
)

type EmailRequest struct {
	Email       string  `json:"email"`
	Type        string  `json:"type"`
	ProjectName string  `json:"project_name"`
	Amount      float64 `json:"amount"`
}
