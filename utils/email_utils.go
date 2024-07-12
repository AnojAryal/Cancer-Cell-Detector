package utils

import (
	"fmt"
	"net/smtp"
)

var (
	smtpServer    = "smtp.gmail.com"
	smtpPort      = "587"
	senderEmail   = "teamcancercelldetector@gmail.com"
	emailPassword = "rgck vvtb uefq ensk"
)

func sendEmail(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", senderEmail, emailPassword, smtpServer)

	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpServer, smtpPort),
		auth,
		senderEmail,
		to,
		[]byte(message),
	)

	if err != nil {
		fmt.Printf("Failed to send email to %v: %v\n", to, err)
		return err
	}

	fmt.Printf("Email sent successfully to %v\n", to)
	return nil
}

func SendVerificationEmail(email, token string) error {
	verificationLink := fmt.Sprintf("http://localhost:3000/users/verify/%s", token)
	subject := "Verify your email"
	body := fmt.Sprintf("Hi,\n\nYour account has been registered in our system. Please click the following link to verify your email:\n%s\n\nThank you!", verificationLink)
	return sendEmail([]string{email}, subject, body)
}

func SendPasswordResetEmail(email, token string) error {
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password/%s", token)
	subject := "Password Reset"
	body := fmt.Sprintf("Hi,\n\nYou have requested to reset your password. Please click on the following link to reset your password:\n%s\n\nIf you didn't request this, you can safely ignore this email.\n\nThank you!", resetLink)
	return sendEmail([]string{email}, subject, body)
}
