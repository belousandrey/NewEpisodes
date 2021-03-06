package matchdaybiz

import (
	"io"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	app "github.com/belousandrey/new-episodes"
	"github.com/pkg/errors"
)

const (
	// MatchDayDateFormat - date format that used at matchday.biz
	MatchDayDateFormat = "02.01.2006 в 15:04"
)

// Engine - special engine for podcast matchday.biz
type Engine struct {
	// lastEpisode - last episode that we know
	LastEpisode string
}

// NewEngine - create new engine
func NewEngine(last string) app.Episoder {
	return &Engine{
		LastEpisode: last,
	}
}

// GetNewEpisodes - find new episodes since LastEpisode
func (e *Engine) GetNewEpisodes(resp io.Reader) (episodes []app.Episode, last string, err error) {
	// parse date from specific date format
	tle, err := time.Parse(app.DateFormat, e.LastEpisode)
	if err != nil {
		err = errors.Wrap(err, "parse date from string")
		return
	}

	// parse HTML document content
	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		err = errors.Wrap(err, "parse HTML document")
		return
	}

	// extract data from table
	doc.Find(".myform .myhistory").EachWithBreak(func(i int, table *goquery.Selection) bool {
		table.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
			// first row has only headers
			if i == 0 {
				return true
			}

			var (
				title, link string
				date        time.Time
			)

			tr.Children().Each(func(i int, td *goquery.Selection) {
				if i == 2 {
					// unformatted string with date
					date, err = time.Parse(MatchDayDateFormat, td.Text())
					if err != nil {
						return
					}
				} else if i == 3 {
					// episode title
					a := td.Find("a")
					title = strings.Replace(strings.TrimSpace(a.Text()), "\t", "", -1)
				} else if i == 7 {
					// link to file
					link, _ = td.Find("a").Attr("href")
				}
			})

			// search for new episodes
			if date.Before(tle) {
				return false
			}

			if last == "" {
				last = date.Add(time.Hour * 24).Format(app.DateFormat)
			}

			ep := app.Episode{
				Link:  link,
				Title: title,
				Date:  date.Format(app.DateFormat),
			}
			episodes = append([]app.Episode{ep}, episodes...)
			return true
		})

		if last != "" {
			return false
		}

		return true
	})

	return
}
