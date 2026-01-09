package internal

import (
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

func SendResetEmail(to, token string) error {
	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		atoi(os.Getenv("SMTP_PORT")),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Password Reset Request")
	resetURL := os.Getenv("RESET_URL") + "?token=" + token
	m.SetBody("text/html", "Click <a href='"+resetURL+"'>here</a> to reset your password.")
	return d.DialAndSend(m)
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
