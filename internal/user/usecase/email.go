package usecase

import (
	"net/smtp"
	"os"
)

type SMTPEmailSender struct {
	From     string
	Password string
	Host     string
	Port     string
}

// ✅ implement EmailSender interface
func (s *SMTPEmailSender) Send(to, subject, body string) error {
	addr := s.Host + ":" + s.Port
	auth := smtp.PlainAuth("", s.From, s.Password, s.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
}

// ✅ สร้าง constructor แบบโหลดจาก .env
func NewEmailSender() EmailSender {
	return &SMTPEmailSender{
		From:     os.Getenv("SMTP_FROM"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
	}
}
