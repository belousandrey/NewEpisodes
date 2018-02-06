package matchdaybiz

import (
	"net/http"

	"strings"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/belousandrey/NewEpisodes/src/const"
	"github.com/belousandrey/NewEpisodes/src/types"
	"github.com/pkg/errors"
)

const (
	MatchDayDateFormat = "02.01.2006 в 15:04"
)

type Engine struct {
	lastEpisode string
}

func NewEngine(last string) *Engine {
	return &Engine{
		lastEpisode: last,
	}
}

func (e *Engine) GetNewEpisodes(resp *http.Response) (episodes []types.Episode, last string, err error) {
	tle, err := time.Parse(constants.DateFormat, e.lastEpisode)
	if err != nil {
		err = errors.Wrap(err, "parse date from string")
		return
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		err = errors.Wrap(err, "parse HTML document")
		return
	}

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

			if date.Before(tle) {
				return false
			}

			if last == "" {
				last = date.Add(time.Hour * 24).Format(constants.DateFormat)
			}

			ep := types.Episode{
				Link:  link,
				Title: title,
				Date:  date.Format(constants.DateFormat),
			}
			episodes = append([]types.Episode{ep}, episodes...)
			return true
		})

		if last != "" {
			return false
		}

		return true
	})

	return
}
