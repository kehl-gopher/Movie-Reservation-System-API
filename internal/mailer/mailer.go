package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	Dialer *mail.Dialer
	Sender string
}

func NewMailer(host, username, sender, password string, port int) *Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 3 * time.Second

	return &Mailer{Dialer: dialer, Sender: sender}
}

func (m *Mailer) Send(recipent, templateFile string, data interface{}) error {

	temp, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = temp.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = temp.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = temp.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()

	// set headers
	msg.SetHeader("To", recipent)
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.Dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
