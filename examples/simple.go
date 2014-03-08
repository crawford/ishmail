package main

import (
	"github.com/crawford/ishmail"
	"html/template"
	"log"
	"net/mail"
	"net/smtp"
)

func main() {
	template, err := template.New("foo").Parse(`{{define "Body"}}<b>{{.}}</b>{{end}}`)
	if err != nil {
		log.Fatal(err)
	}
	msg := ishmail.CreateHtmlEmail(
		"Test",
		"Hello world!",
		template,
		&mail.Address{
			Name:    "Sender",
			Address: "sender@example.com",
		},
		&mail.Address{
			Name:    "Receiver One",
			Address: "receiver1@example.com",
		},
		&mail.Address{
			Name:    "Receiver Two",
			Address: "receiver2@example.com",
		},
	)

	auth := smtp.PlainAuth("<identity>",
		"<username>",
		"<password>",
		"<host>")
	err = ishmail.Send(msg, auth, "<remote addr>:<remote port>")
	if err != nil {
		log.Fatal(err)
	}
}
