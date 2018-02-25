package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"strconv"

	"runtime"

	"path"

	"github.com/belousandrey/new-episodes/src/types"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
)

const emailTemplateFile = "../templates/email.html"

// SendEmail - generate HTML from template, send email
func SendEmail(recipient string, sender map[string]string, content *types.EmailContent) error {
	// process template
	html, err := PrepareTemplate(content)
	if err != nil {
		return errors.Wrap(err, "prepare template")
	}

	port, err := strconv.Atoi(sender["port"])
	if err != nil {
		return errors.Wrap(err, "convert string to integer")
	}
	d := gomail.NewDialer(sender["smtp"], port, sender["username"], sender["password"])

	// send email
	m := gomail.NewMessage()
	m.SetHeader("From", sender["username"]+"@"+sender["domain"])
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "New podcast episodes")
	m.SetBody("text/html", html.String())

	if err := d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "dial and send email")
	}
	return nil
}

// PrepareTemplate - parse template file, fill it with data
func PrepareTemplate(content *types.EmailContent) (*bytes.Buffer, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return nil, errors.New("get program location")
	}

	t, err := template.ParseFiles(path.Join(path.Dir(filename), emailTemplateFile))
	if err != nil {
		return nil, errors.Wrap(err, "parse template file")
	}

	var html bytes.Buffer
	if err = t.Execute(&html, content); err != nil {
		return nil, errors.Wrap(err, "execute template content")
	}

	return &html, nil
}

// DownloadFile - download file by provided URL
func DownloadFile(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create request object")
	}

	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "download file by URL")
	}

	return resp.Body, nil
}
