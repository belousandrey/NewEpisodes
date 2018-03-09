package defaultengine

import (
	"io"
	"time"

	app "github.com/belousandrey/new-episodes"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// Engine - default engine for most podcasts
type Engine struct {
	// lastEpisode - last episode that we know
	LastEpisode string
}

// NewEngine - create new default engine
func NewEngine(last string) app.Episoder {
	return &Engine{
		LastEpisode: last,
	}
}

// GetNewEpisodes - find new episodes since LastEpisode
func (e *Engine) GetNewEpisodes(resp io.Reader) (episodes []app.Episode, last string, err error) {
	// parse date from default date format
	tle, err := time.Parse(app.DateFormat, e.LastEpisode)
	if err != nil {
		err = errors.Wrap(err, "parse date from string")
		return
	}

	// parse RSS content
	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp)
	if err != nil {
		err = errors.Wrap(err, "parse RSS feed body")
		return
	}

	// search for new episodes
	for _, e := range feed.Items {
		if e.PublishedParsed.Before(tle) {
			break
		}

		if last == "" {
			last = e.PublishedParsed.Add(time.Hour * 24).Format(app.DateFormat)
		}

		ep := app.Episode{
			Title: e.Title,
			Link:  e.Link,
			Date:  e.PublishedParsed.Format(app.DateFormat),
		}
		episodes = append([]app.Episode{ep}, episodes...)
	}

	return
}
