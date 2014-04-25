/*
Copyright 2014 Alex Crawford

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ishmail

import (
	"errors"
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
	message, err := generateMessage(email)
	if err != nil {
		return err
	}
	return smtp.SendMail(addr, auth, email.From().Address, addrsToAddresses(email.To()), message)
}

func generateMessage(email Emailer) ([]byte, error) {
	if email == nil {
		return nil, errors.New("Email must not be nil")
	}

	if email.From() == nil {
		return nil, errors.New("From must not be nil")
	}

	if email.To() == nil {
		return nil, errors.New("To must not be nil")
	}

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
		return nil, err
	}
	return []byte(headersToString(headers) + body), nil
}

func headersToString(headers map[string]string) string {
	str := ""
	for k, v := range headers {
		str += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return str + "\r\n"
}

func addrsToAddresses(addrs []*mail.Address) []string {
	out := make([]string, 0)
	for _, addr := range addrs {
		out = append(out, addr.Address)
	}
	return out
}

func addrsToString(addrs []*mail.Address) string {
	out := ""
	if len(addrs) > 0 {
		out += addrs[0].String()
		addrs = addrs[1:]
	}
	for _, addr := range addrs {
		out += ", " + addr.String()
	}
	return out
}
