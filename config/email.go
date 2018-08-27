package config

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

// SendEmail using gmail server nmtp
func SendEmail(to string, subject string, tem string) error {
	c := GetConfig()

	from := mail.Address{c.Email.Name, c.Email.From}
	toMail := mail.Address{"", to}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = toMail.String()
	headers["Subject"] = subject
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += tem

	// authentication
	auth := smtp.PlainAuth("", c.Email.From, c.Email.Password, c.Email.Host)

	// Config tls security
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.Email.Host,
	}
	conn, err := tls.Dial("tcp", c.Email.Server, tlsConfig)
	if err != nil {
		return err
	}

	// ---------------------------------------------------------------------------
	//  create new client
	client, err := smtp.NewClient(conn, c.Email.Host)
	if err != nil {
		return err
	}

	// Authenticate
	err = client.Auth(auth)
	if err != nil {
		return err
	}

	// Set From email
	err = client.Mail(from.Address)
	if err != nil {
		return err
	}

	// Set to email recipient
	err = client.Rcpt(toMail.Address)
	if err != nil {
		return err
	}

	// Process data send
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	// Close client
	err = w.Close()
	if err != nil {
		return err
	}

	// Exit client
	client.Quit()

	return nil
}
