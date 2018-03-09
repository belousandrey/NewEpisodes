package refresher

import (
	"bytes"
	"html/template"
	"path"
	"runtime"
	"strconv"

	app "github.com/belousandrey/new-episodes"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
)

const emailTemplateFile = "../templates/email.html"

// SendEmail - generate HTML from template, send email
func SendEmail(recipient string, sender map[string]string, content *app.EmailContent) error {
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
func PrepareTemplate(content *app.EmailContent) (*bytes.Buffer, error) {
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
