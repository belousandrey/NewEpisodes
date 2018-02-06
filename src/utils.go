package main

import (
	"bytes"
	"html/template"
	"net/http"

	"strconv"

	"runtime"

	"path"

	"github.com/belousandrey/NewEpisodes/src/types"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
)

var emailTemplateFile = "../templates/email.html"

func sendEmail(recepient string, sender map[string]string, data []types.PodcastWithEpisodes) error {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("get program location")
	}

	t, err := template.ParseFiles(path.Join(path.Dir(filename), emailTemplateFile))
	if err != nil {
		return errors.Wrap(err, "parse template file")
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
