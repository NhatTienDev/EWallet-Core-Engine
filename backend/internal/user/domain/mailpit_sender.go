package domain

type MailpitSenderRepository interface {
	SendResetPasswordEmail(toEmail, resetToken string) error
	SendPasswordChangedAlert(toEmail string) error
}