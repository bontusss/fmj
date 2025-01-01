package email

import (
	"fmj/config"
	"fmt"
	"gopkg.in/mail.v2"
	"log"
)

type Service interface {
	SendVerificationEmail(to, name, code string) error
	SendWelcomeEmail(to, name string) error
}

type service struct {
	config *config.Config
}

func (s *service) SendVerificationEmail(to, name, code string) error {
	subject := "Verify your email address"
	verifyLink := fmt.Sprintf("%s/auth/verify?code=%s", s.config.BaseURL, code)
	body := fmt.Sprintf("Hello %s,\n\nPlease verify your email by clicking this link: %s", name, verifyLink)

	return s.sendEmail(to, subject, body)
}

func (s *service) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to our platform!"
	body := fmt.Sprintf("Hello %s,\n\nWelcome to our platform. We're excited to have you!", name)

	return s.sendEmail(to, subject, body)
}

func (s *service) sendEmail(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.config.FromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := mail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUsername, s.config.SMTPPassword)

	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	fmt.Println("Email sent")
	return nil
}

func NewService(config *config.Config) Service {
	return &service{config: config}
}
