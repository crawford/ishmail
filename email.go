package ishmail

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"time"
)

// Partial implementation of the emailer interface.
type Email struct {
	from *mail.Address
	to   []*mail.Address
}

func (e *Email) From() *mail.Address {
	return e.from
}

func (e *Email) To() []*mail.Address {
	return e.to
}

func (e *Email) Date() time.Time {
	return time.Now()
}

// An emailer is able to generate the email's body and all of its associated
// metadata (sender, receivers, subject, etc.).
type Emailer interface {
	From() *mail.Address
	To() []*mail.Address
	Subject() string
	ContentType() string
	Body() (string, error)
	Date() time.Time
}

// Send an email with the specified authentication to a host.
// This will not properly escape all of the fields or the message body.
func Send(email Emailer, auth smtp.Auth, addr string) error {
	headers := map[string]string{
		"From":         email.From().String(),
		"To":           addrsToString(email.To()),
		"Subject":      email.Subject(),
		"MIME-Version": "1.0",
		"Content-Type": email.ContentType(),
		"Date":         email.Date().Format(time.RFC1123Z),
	}

	body, err := email.Body()
	if err != nil {
		return err
	}
	message := []byte(headersToString(headers) + body)

	return smtp.SendMail(addr, auth, email.From().Address, addrsToAddress(email.To()), message)
}

func headersToString(headers map[string]string) string {
	str := ""
	for k, v := range headers {
		str += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return str
}

func addrsToAddress(addrs []*mail.Address) []string {
	out := make([]string, 0)
	for _, addr := range addrs {
		out = append(out, addr.Address)
	}
	return out
}

func addrsToString(addrs []*mail.Address) string {
	out := ""
	for _, addr := range addrs {
		out += addr.String() + ", "
	}
	return out
}
