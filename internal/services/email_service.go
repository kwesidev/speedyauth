package services

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/kwesidev/speedyauth/internal/models"
	"gopkg.in/gomail.v2"
)

type EmailServiceInterface interface {
	SendTwoFactorRequest(randomCodes string, userDetails models.User) error
	SendEmailLoginRequest(randomCodes string, userDetails models.User) error
	SendPasswordResetRequest(randomCodes string, userDetails models.User) error
}
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
func (emSrv *EmailService) sendEmail(to []string, subject, message string) error {
	portNumber, _ := strconv.Atoi(emSrv.smtpPort)
	d := gomail.NewDialer(emSrv.smtpHost, portNumber, emSrv.smtpUsername, emSrv.smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: emSrv.secure}
	// Compose the message to be sent
	m := gomail.NewMessage()
	m.SetHeader("From", emSrv.fromEmailAddress)
	m.SetHeader("To", to[:]...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)
	// Proceed to send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// SendTwoFactorRequest sends two factor mail
func (emSrv *EmailService) SendTwoFactorRequest(randomCodes string, userDetails models.User) error {
	var twoFactorRequestTemplateBuffer bytes.Buffer
	// Get email template from directory and assign random code to it
	emailTemplateFile, err := template.ParseFiles("static/email_templates/TwoFactorLogin.html")
	if err != nil {
		return err
	}
	tmpl := template.Must(emailTemplateFile, err)
	emailTemplateData := struct {
		FullName   string
		RandomCode string
	}{}
	emailTemplateData.RandomCode = randomCodes
	emailTemplateData.FullName = userDetails.FirstName + " " + userDetails.LastName
	tmpl.Execute(&twoFactorRequestTemplateBuffer, emailTemplateData)
	recipient := []string{userDetails.EmailAddress}
	if err = emSrv.sendEmail(recipient, "Two-factor login", twoFactorRequestTemplateBuffer.String()); err != nil {
		log.Println("Sending Two Factor Request Email Error", err)
		return err
	}
	return nil
}

// SendTwoFactorRequest sends two factor mail
func (emSrv *EmailService) SendEmailLoginRequest(randomCodes string, userDetails models.User) error {
	var twoFactorRequestTemplateBuffer bytes.Buffer
	// Get email template from directory and assign random code to it
	emailTemplateFile, err := template.ParseFiles("static/email_templates/EmailLogin.html")
	if err != nil {
		return err
	}
	tmpl := template.Must(emailTemplateFile, err)
	emailTemplateData := struct {
		FullName   string
		RandomCode string
	}{}
	emailTemplateData.RandomCode = randomCodes
	emailTemplateData.FullName = userDetails.FirstName + " " + userDetails.LastName
	tmpl.Execute(&twoFactorRequestTemplateBuffer, emailTemplateData)
	recipient := []string{userDetails.EmailAddress}
	if err = emSrv.sendEmail(recipient, "Email login", twoFactorRequestTemplateBuffer.String()); err != nil {
		log.Println("Sending Email Login Request  Error", err)
		return err
	}
	return nil
}

// SendPasswordRequest
// Sends a password request mail to the receiver
func (emSrv *EmailService) SendPasswordResetRequest(randomCodes string, userDetails models.User) error {
	var passwordResetTemplateBuffer bytes.Buffer
	// Get email template from directory and assign random code to it
	emailTemplateFile, err := template.ParseFiles("static/email_templates/PasswordRequest.html")
	if err != nil {
		log.Println("Template reading ", err)
		return err
	}
	tmpl := template.Must(emailTemplateFile, err)
	emailTemplateData := struct {
		FullName   string
		RandomCode string
	}{}
	emailTemplateData.RandomCode = randomCodes
	emailTemplateData.FullName = userDetails.FirstName + " " + userDetails.LastName
	tmpl.Execute(&passwordResetTemplateBuffer, emailTemplateData)
	recipient := []string{userDetails.EmailAddress}
	if err = emSrv.sendEmail(recipient, "Password Reset Request", passwordResetTemplateBuffer.String()); err != nil {
		log.Println("Sending Password Reset Email Error", err)
		return err
	}
	return nil
}
