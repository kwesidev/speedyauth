package services

import (
	"crypto/tls"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	smtpHost         string
	smtpUsername     string
	smtpPassword     string
	smtpPort         string
	fromEmailAddress string
	secure           bool
}

func NewEmailService(secure bool) *EmailService {
	return &EmailService{
		smtpHost:         os.Getenv("SMTP_HOST"),
		smtpUsername:     os.Getenv("SMTP_USERNAME"),
		smtpPort:         os.Getenv("SMTP_PORT"),
		smtpPassword:     os.Getenv("SMTP_PASSWORD"),
		fromEmailAddress: os.Getenv("FROM_EMAIL_ADDRESS"),
		secure:           secure,
	}
}

// SendEmail funnction sends email directly to an external server
func (this *EmailService) SendEmail(to []string, subject, message string) error {
	portNumber, _ := strconv.Atoi(this.smtpPort)
	d := gomail.NewDialer(this.smtpHost, portNumber, this.smtpUsername, this.smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: this.secure}
	// Compose the message to be sent
	m := gomail.NewMessage()
	m.SetHeader("From", this.fromEmailAddress)
	m.SetHeader("To", to[:]...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)
	// Proceed to send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
