package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/belousandrey/NewEpisodes/src/types"
	"github.com/go-gomail/gomail"
)

var emailTemplate = `
New episodes in these podcasts:

{{range .}}{{.Podcast.Title}}{{range .Episodes}}<ul>
<li><a href="{{.Link}}">{{.Title}}</a> at {{.Date}}</li></ul>{{end}}
{{end}}
`

func sendEmail(address string, data []*types.PodcastWithEpisodes) error {
	t := template.New("main")
	t, err := t.Parse(emailTemplate)
	if err != nil {
		return err
	}

	var html bytes.Buffer
	if err = t.Execute(&html, data); err != nil {
		return err
	}

	fmt.Printf("send this to address %s\n", address)
	fmt.Println(html.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, "someguy", "password")

	m := gomail.NewMessage()
	m.SetHeader("From", address)
	m.SetHeader("To", address)
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
