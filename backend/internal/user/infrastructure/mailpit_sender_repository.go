package infrastructure

import (
	"fmt"
	"net/smtp"

	"github.com/nhattiendev/ewallet/internal/user/domain"
)

type mailpitSenderRepository struct {
	smtpHost string
	smtpPort string
}

func NewMailpitSenderRepository(smtpHost, smtpPort string) domain.MailpitSenderRepository {
	return &mailpitSenderRepository{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (m *mailpitSenderRepository) SendResetPasswordEmail(toEmail, resetToken string) error {
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)
	
	subject := "Subject: EWallet Password Reset Request\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<h2>Reset Your Password</h2>
		<p>You have requested to reset your password. Please click the link below (valid for 15 minutes):</p>
		<a href="%s">Reset Password</a>
		<p>If you did not request this, please ignore this email. Your password will remain unchanged.</p>
	`, resetLink)

	msg := []byte(subject + mime + body)

	// Mailpit does not require authentication, can pass auth = nil
	addr := fmt.Sprintf("%s:%s", m.smtpHost, m.smtpPort)
	return smtp.SendMail(addr, nil, "noreply@ewallet.local", []string{toEmail}, msg)
}

func (m *mailpitSenderRepository) SendPasswordChangedAlert(toEmail string) error {
	subject := "Subject: Security Alert: Your Password Has Been Changed\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := `
		<h2>Security Alert</h2>
		<p>Your EWallet account password has been successfully changed.</p>
		<p>If you did not make this change, please contact our support team immediately to secure your account!</p>
	`

	msg := []byte(subject + mime + body)
	addr := fmt.Sprintf("%s:%s", m.smtpHost, m.smtpPort)
	return smtp.SendMail(addr, nil, "security@ewallet.local", []string{toEmail}, msg)
}