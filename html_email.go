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
	"html/template"
	"net/mail"
)

// HtmlEmail provides a basic type for creating HTML-based emails. The
// provided content is fed into the template in order to generate the body of
// the message.
type HtmlEmail struct {
	subject  string
	content  interface{}
	template *template.Template
	Email
}

func CreateHtmlEmail(subject string, content interface{}, template *template.Template, from *mail.Address, to ...*mail.Address) *HtmlEmail {
	email := &HtmlEmail{
		subject:  subject,
		content:  content,
		template: template,
	}
	email.from = from
	email.to = to
	return email
}

// Generate the body of the message from the given content and template
func (e *HtmlEmail) Body() (string, error) {
	var buffer bytes.Buffer
	err := e.template.ExecuteTemplate(&buffer, "Body", e.content)
	return buffer.String(), err
}

// Returns "text/html"
func (e *HtmlEmail) ContentType() string {
	return "text/html"
}

// Returns the given subject
func (e *HtmlEmail) Subject() string {
	return e.subject
}
