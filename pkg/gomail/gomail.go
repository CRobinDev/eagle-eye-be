package gomail

import (
	"bytes"
	templ "html/template"
	"os"
	"path/filepath"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/template"
	"gopkg.in/gomail.v2"
)

type Gomail struct {
	message   *gomail.Message
	dialer    *gomail.Dialer
	templates *templ.Template
}

func NewGomail() *Gomail {
	templates, err := templ.ParseFS(template.Template, "*.html")
	if err != nil {
		return &Gomail{}
	}

	return &Gomail{
		message:   gomail.NewMessage(),
		dialer:    gomail.NewDialer(env.GetEnv().SmtpHost, env.GetEnv().SmtpPort, env.GetEnv().SmtpEmail, env.GetEnv().SmtpPassword),
		templates: templates,
	}
}

func (g *Gomail) SetBodyHTML(path string, data interface{}) (string, error) {
	var body bytes.Buffer

	t := g.templates.Lookup(path)
	if t == nil {
		return "", errorz.ErrSetHTMLTemplate
	}

	err := t.Execute(&body, data)
	if err != nil {
		return "", errorz.ErrExecuteHTML
	}

	return body.String(), nil
}

func (g *Gomail) Send(custEmail, subject, path string, data map[string]interface{}) error {
	emailBody, err := g.SetBodyHTML(path, data)
	if err != nil {
		return err
	}

	g.message.SetHeader("From", "Eagle Eye Team <"+env.GetEnv().SmtpEmail+">")
	g.message.SetHeader("To", custEmail)
	g.message.SetHeader("Subject", subject)
	g.message.SetBody("text/html", emailBody)

	tmpPath := filepath.Join(os.TempDir(), "logo.png")
	if err := os.WriteFile(tmpPath, template.Logo, 0644); err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	g.message.Embed(tmpPath)

	if err := g.dialer.DialAndSend(g.message); err != nil {
		return errorz.ErrFailedToSendNotification
	}

	return nil
}
