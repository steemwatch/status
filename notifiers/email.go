package notifiers

import (
	"crypto/tls"
	"fmt"

	"github.com/steemwatch/status/checks"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

//
// Config
//

type EmailNotifierConfig struct {
	SMTPServerHost string   `yaml:"smtp_server_host"`
	SMTPServerPort int      `yaml:"smtp_server_port"`
	SMTPUsername   string   `yaml:"smtp_username"`
	SMTPPassword   string   `yaml:"smtp_password"`
	From           string   `yaml:"from"`
	To             []string `yaml:"to"`
}

func (config *EmailNotifierConfig) Validate() error {
	switch {
	case config.SMTPServerHost == "":
		return errors.New("key not set: email.smtp_server_host")
	case config.SMTPServerPort == 0:
		return errors.New("key not set: email.smtp_server_port")
	case config.SMTPUsername == "":
		return errors.New("key not set: email.smtp_username")
	case config.SMTPPassword == "":
		return errors.New("key not set: email.smtp_password")
	case config.From == "":
		return errors.New("key not set: email.from")
	case len(config.To) == 0:
		return errors.New("array empty: email.to")
	default:
		return nil
	}
}

//
// Notifier
//

type EmailNotifier struct {
	config *EmailNotifierConfig
}

func NewEmailNotifier(config *EmailNotifierConfig) (*EmailNotifier, error) {
	// Validate.
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Done.
	return &EmailNotifier{config}, nil
}

func (notifier *EmailNotifier) DispatchNotification(check *checks.CheckSummary) error {
	subject := "[SteemWatch] Check changed to " + string(check.Result)
	body := fmt.Sprintf(`
Description: %v
Result:      %v
Details:     %v
Timestamp:   %v
`, check.Description, check.Result, check.Details, check.Timestamp)

	return notifier.send(subject, body, "text/plain")
}

func (notifier *EmailNotifier) send(
	subject string,
	body string,
	contentType string,
) error {

	config := notifier.config

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.From)
	msg.SetHeader("To", config.To...)
	msg.SetHeader("Subject", subject)
	msg.SetBody(contentType, body)

	dialer := gomail.NewDialer(
		config.SMTPServerHost, config.SMTPServerPort,
		config.SMTPUsername, config.SMTPPassword)

	dialer.TLSConfig = &tls.Config{ServerName: config.SMTPServerHost}
	return dialer.DialAndSend(msg)
}
