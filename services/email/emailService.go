package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"

	"bbscout/types"
)

func SendEmail(email types.EmailPayload) {
	godotenv.Load()
	tmpl, err := template.ParseFiles("template/" + email.TemplateFile)
	if err != nil {
		fmt.Println("failed to get template")
		return
	}

	// Create a buffer to hold the rendered template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, email); err != nil {
		fmt.Println("could not execute template")
	}

	// Set up the email
	m := gomail.NewMessage()
	m.SetAddressHeader("From", getEnv("MAIL_FROM", "dunstansafu@gmail.com"), "BBscout")
	m.SetHeader("To", email.MailTo)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", body.String())

	// Attach the image and set its Content-ID
	// m.Embed("template/diracks-header.png", gomail.SetHeader(map[string][]string{
	// 	"Content-ID":          {"<welcomeImage>"},
	// 	"Content-Disposition": {"inline"},
	// }))

	for _, attachment := range email.Attatchments {
		m.Attach(attachment)
	}

	port, _ := strconv.Atoi(getEnv("MAIL_PORT", "465"))

	// SMTP server configuration
	d := gomail.NewDialer(getEnv("MAIL_SERVER", "smtp.fastmail.com"), port, getEnv("MAIL_USERNAME", "dunstansafu@gmail.com"), getEnv("MAIL_PASSWORD", "kjx7r6n7vzqwsmtz"))
	d.SSL = true

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("could not send email %s", err)
	}

}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
