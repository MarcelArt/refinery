package repositories

import (
	"fmt"
	"strings"

	"git.bangmarcel.art/marcel/arrays"
	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"gopkg.in/gomail.v2"
)

type MailRepo struct {
	dialer   *gomail.Dialer
	smtpName string
}

func NewMailRepo() *MailRepo {
	dialer := gomail.NewDialer(
		configs.Env.SMTPHost,
		configs.Env.SMTPPort,
		configs.Env.SMTPEmail,
		configs.Env.SMTPPassword,
	)

	return &MailRepo{
		dialer:   dialer,
		smtpName: configs.Env.SMTPName,
	}
}

func (r *MailRepo) SendMail(m models.Mailer) error {
	if configs.Env.ServerENV != "prod" {
		m.To = arrays.Map(m.To, func(to string) string {
			emailName, _, _ := strings.Cut(to, "@")
			return fmt.Sprintf("%s@yopmail.com", emailName)
		})
	}

	message := gomail.NewMessage()
	message.SetHeader("From", r.smtpName)
	message.SetHeader("To", m.To...)
	message.SetHeader("Subject", m.Subject)
	message.SetBody("text/html", m.Body)

	for _, attachment := range m.Attachments {
		message.Attach(attachment)
	}

	return r.dialer.DialAndSend(message)
}
