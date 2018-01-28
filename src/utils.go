package main

import (
	"bytes"
	"html/template"
	"net/http"

	"strconv"

	"github.com/belousandrey/NewEpisodes/src/types"
	"github.com/go-gomail/gomail"
)

var emailTemplateFile = "../templates/email.html"

func sendEmail(recepient string, sender map[string]string, data []*types.PodcastWithEpisodes) error {
	t, err := template.ParseFiles(emailTemplateFile)
	if err != nil {
		return err
	}

	var html bytes.Buffer
	if err = t.Execute(&html, data); err != nil {
		return err
	}

	port, err := strconv.Atoi(sender["port"])
	if err != nil {
		return err
	}
	d := gomail.NewDialer(sender["host"], port, sender["username"], sender["password"])

	m := gomail.NewMessage()
	m.SetHeader("From", recepient)
	m.SetHeader("To", recepient)
	m.SetHeader("Subject", "New podcast episodes")
	m.SetBody("text/html", html.String())

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil
}

func DownloadFile(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
