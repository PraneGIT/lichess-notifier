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
    
    message.SetHeader("From", e.from)
    message.SetHeader("To", e.to)
    message.SetHeader("Subject", subject)
    
    message.SetBody("text/plain", body)

    dialer := gomail.NewDialer(e.smtpHost, e.smtpPort, e.from, e.password)
    return dialer.DialAndSend(message)
}