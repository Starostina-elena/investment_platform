package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/smtp"
	"os"

	"github.com/Starostina-elena/investment_platform/services/notification/core"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	log          slog.Logger
}

func NewEmailService(log slog.Logger) *EmailService {
	return &EmailService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     os.Getenv("SMTP_PORT"),
		smtpUser:     os.Getenv("SMTP_USER"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromEmail:    os.Getenv("FROM_EMAIL"),
		log:          log,
	}
}

func (s *EmailService) SendNotification(req *core.EmailRequest) error {
	subject, body, err := s.buildEmail(req)
	if err != nil {
		s.log.Error("failed to build email", "error", err)
		return err
	}

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
		"%s", s.fromEmail, req.Email, subject, body))

	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)
	addr := s.smtpHost + ":" + s.smtpPort

	err = smtp.SendMail(addr, auth, s.fromEmail, []string{req.Email}, msg)
	if err != nil {
		s.log.Error("failed to send email", "error", err, "to", req.Email)
		return core.ErrEmailSendFailed
	}

	s.log.Info("email sent", "to", req.Email, "type", req.Type)
	return nil
}

func (s *EmailService) buildEmail(req *core.EmailRequest) (string, string, error) {
	switch req.Type {
	case core.NotifTypeDividends:
		return s.buildDividendsEmail(req)
	case core.NotifTypeRefund:
		return s.buildRefundEmail(req)
	default:
		return "", "", core.ErrUnknownNotifType
	}
}

func (s *EmailService) buildDividendsEmail(req *core.EmailRequest) (string, string, error) {
	subject := "Выплата дивидендов"
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #4CAF50;">Выплата дивидендов</h2>
        <p>Здравствуйте!</p>
        <p>Вам выплачены дивиденды по проекту <strong>{{.ProjectName}}</strong>.</p>
        <p style="font-size: 18px; color: #4CAF50;">
            <strong>Сумма: {{.Amount}} ₽</strong>
        </p>
        <p>Средства зачислены на ваш счет.</p>
        <hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
        <p style="font-size: 12px; color: #888;">
            Это автоматическое уведомление, не отвечайте на него.
        </p>
    </div>
</body>
</html>
`
	t, err := template.New("dividends").Parse(tmpl)
	if err != nil {
		return "", "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, req); err != nil {
		return "", "", err
	}

	return subject, buf.String(), nil
}

func (s *EmailService) buildRefundEmail(req *core.EmailRequest) (string, string, error) {
	subject := "Возврат средств"
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2196F3;">Возврат средств</h2>
        <p>Здравствуйте!</p>
        <p>Проект <strong>{{.ProjectName}}</strong> завершен.</p>
        <p style="font-size: 18px; color: #2196F3;">
            <strong>Возвращено: {{.Amount}} ₽</strong>
        </p>
        <p>Ваши инвестиционные средства возвращены на счет.</p>
        <hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
        <p style="font-size: 12px; color: #888;">
            Это автоматическое уведомление, не отвечайте на него.
        </p>
    </div>
</body>
</html>
`
	t, err := template.New("refund").Parse(tmpl)
	if err != nil {
		return "", "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, req); err != nil {
		return "", "", err
	}

	return subject, buf.String(), nil
}
