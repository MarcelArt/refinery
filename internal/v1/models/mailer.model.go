package models

type Mailer struct {
	To          []string
	Subject     string
	Body        string
	Attachments []string
}
