package notifier

import (
	"github.com/PraneGIT/lichess-notifier/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	from     string
	to       string
	password string
	smtpHost string
	smtpPort int
}

func NewEmailNotifier(thisConfig config.EmailConfig) *EmailNotifier {
	return &EmailNotifier{
		from:     thisConfig.From,
		to:       thisConfig.To,
		password: thisConfig.Password,
		smtpHost: thisConfig.SMTPHost,
		smtpPort: thisConfig.SMTPPort,
	}
}

func (e *EmailNotifier) SendEmail(subject string, body string) error {
	message := gomail.NewMessage()
	
	message.SetBody("From",e.from)
	message.SetBody("To",e.to)
	message.SetBody("Text/plain",body)
	message.SetBody("subject",subject)

	dialer := gomail.NewDialer(e.smtpHost,e.smtpPort,e.from,e.to)
	return dialer.DialAndSend(message)
}
