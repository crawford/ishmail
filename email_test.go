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
	"bytes"
	"errors"
	"io/ioutil"
	"net/mail"
	"testing"
	"time"
)

type TestEmail struct {
	from        *mail.Address
	to          []*mail.Address
	subject     string
	contentType string
	body        string
	bodyErr     error
	date        time.Time
}

func (e TestEmail) From() *mail.Address {
	return e.from
}

func (e TestEmail) To() []*mail.Address {
	return e.to
}

func (e TestEmail) Subject() string {
	return e.subject
}

func (e TestEmail) ContentType() string {
	return e.contentType
}

func (e TestEmail) Body() (string, error) {
	return e.body, e.bodyErr
}

func (e TestEmail) Date() time.Time {
	return e.date
}

func TestGenerateEmailNil(t *testing.T) {
	if _, err := generateMessage(nil); err == nil {
		t.FailNow()
	}
}

func TestGenerateEmailNilFrom(t *testing.T) {
	email := TestEmail{
		from: nil,
		to:   []*mail.Address{},
	}

	_, err := generateMessage(email)
	if err == nil {
		t.FailNow()
	}
}

func TestGenerateEmailNilTo(t *testing.T) {
	email := TestEmail{
		from: &mail.Address{},
		to:   nil,
	}

	_, err := generateMessage(email)
	if err == nil {
		t.FailNow()
	}
}

func TestGenerateEmailBodyErr(t *testing.T) {
	bodyErr := errors.New("Test")
	email := TestEmail{
		from:    &mail.Address{},
		to:      []*mail.Address{},
		bodyErr: bodyErr,
	}

	_, err := generateMessage(email)
	if err != bodyErr {
		t.FailNow()
	}
}

func TestGenerateEmail(t *testing.T) {
	email := TestEmail{
		from: &mail.Address{
			Name:    "Test One",
			Address: "test1@example.com",
		},
		to: []*mail.Address{
			{
				Name:    "Test Two",
				Address: "test2@example.com",
			},
			{
				Name:    "Test Three",
				Address: "test3@example.com",
			},
		},
		subject:     "Test Email",
		contentType: "text/html",
		body:        "Hello.",
		date:        time.Time{},
	}
	expect := map[string][]string{
		"From":         []string{email.from.String()},
		"To":           []string{email.to[0].String() + ", " + email.to[1].String()},
		"Subject":      []string{"Test Email"},
		"Mime-Version": []string{"1.0"},
		"Content-Type": []string{"text/html"},
		"Date":         []string{"Mon, 01 Jan 0001 00:00:00 +0000"},
	}

	rawMessage, err := generateMessage(email)
	if err != nil {
		t.Fatalf("Failed to generate message (%s)", err)
	}

	message, err := mail.ReadMessage(bytes.NewBuffer(rawMessage))
	if err != nil {
		t.Fatalf("Failed to read message (%s)", err)
	}

	body, err := ioutil.ReadAll(message.Body)
	if err != nil {
		t.FailNow()
	}
	if !bytes.Equal(body, []byte(email.body)) {
		t.Fatal("Body is wrong")
	}

	if len(expect) != len(message.Header) {
		t.Fatal("Header length is wrong")
	}
	for k, v := range expect {
		if len(v) != len(message.Header[k]) {
			t.Fatalf("Header entry is wrong (%s:%s)", k, v)
		}
		for i, s := range v {
			if s != message.Header[k][i] {
				t.Fatalf("Header entry is wrong (%s:%s)", k, v, message.Header[k][i])
			}
		}
	}
}

func TestHeadersToStringNil(t *testing.T) {
	if headersToString(nil) != "\r\n" {
		t.FailNow()
	}
}

func TestHeadersToString(t *testing.T) {
	headers := map[string]string{
		"From":         "test1@example.com <Test Person>",
		"To":           "test2@example.com",
		"Subject":      "Testing",
		"MIME-Version": "1.0",
		"Content-Type": "text/html",
		"Date":         time.Date(2014, time.March, 29, 2, 32, 17, 1232, time.UTC).Format(time.RFC1123Z),
	}
	expect := "From: test1@example.com <Test Person>\r\nTo: test2@example.com\r\nSubject: Testing\r\nMIME-Version: 1.0\r\nContent-Type: text/html\r\nDate: Sat, 29 Mar 2014 02:32:17 +0000\r\n\r\n"
	if headersToString(headers) != expect {
		t.FailNow()
	}
}

func TestAddrsToAddresses(t *testing.T) {
	tests := []struct {
		input  []*mail.Address
		output []string
	}{
		{
			input:  []*mail.Address{},
			output: []string{},
		},
		{
			input: []*mail.Address{
				&mail.Address{
					Name:    "First Last",
					Address: "test1@example.com",
				},
			},
			output: []string{
				"test1@example.com",
			},
		},
		{
			input: []*mail.Address{
				&mail.Address{
					Name:    "First Last",
					Address: "test1@example.com",
				},
				&mail.Address{
					Name:    "A B",
					Address: "test2@example.com",
				},
			},
			output: []string{
				"test1@example.com",
				"test2@example.com",
			},
		},
	}
	for _, test := range tests {
		got := addrsToAddresses(test.input)
		if len(got) != len(test.output) {
			t.Errorf("Expected %d addresses but got %d addresses", len(test.output), len(got))
		}
		for i, result := range got {
			if result != test.output[i] {
				t.Errorf("Expected '%s' but got '%s'", test.output[i], result)
			}
		}
	}
}
