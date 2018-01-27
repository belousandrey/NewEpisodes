package golangshow

import (
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/belousandrey/NewEpisodes/src/const"
	"github.com/belousandrey/NewEpisodes/src/types"
)

const (
	GolangShowDateFormat = "Mon, 02 Jan 2006 00:00:00 +0000"
)

type Engine struct {
	lastEpisode string
}

func NewEngine(last string) *Engine {
	return &Engine{
		lastEpisode: last,
	}
}

func (e *Engine) GetNewEpisodes(resp *http.Response) (episodes []*types.Episode, last string, err error) {
	tle, err := time.Parse(constants.DateFormat, e.lastEpisode)
	if err != nil {
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		return
	}

	for _, e := range feed.Items {
		t, err2 := time.Parse(GolangShowDateFormat, e.Published)
		if err2 != nil {
			return episodes, last, err2
		}

		if t.Before(tle) {
			break
		}

		if last == "" {
			last = t.Add(time.Hour * 24).Format(constants.DateFormat)
		}

		episodes = append([]*types.Episode{types.NewEpisode(e.Title, e.Link, t.Format(constants.DateFormat))}, episodes...)
	}

	return
}

/*func (e *Engine) GetNewEpisodes_(resp *http.Response) (episodes []*types.Episode, last string, err error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}

	doc.Find(".posts .post").EachWithBreak(func(i int, s *goquery.Selection) bool {
		a := s.Find(".post-title a")
		link, exists := a.Attr("href")
		if !exists {
			return true
		}

		parsedLink := parseGolangShowLink(link)
		if parsedLink == e.lastEpisode {
			return false
		}
		if last == "" {
			last = parsedLink
		}

		episodes = append([]*types.Episode{types.NewEpisode(a.Text(), link, "")}, episodes...)
		return true
	})

	return
}

func parseGolangShowLink(url string) string {
	subStrings := strings.Split(url, "/")
	return subStrings[4] + "/" + subStrings[5]
}*/
