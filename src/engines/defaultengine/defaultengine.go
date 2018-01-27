package defaultengine

import (
	"time"

	"net/http"

	"github.com/belousandrey/NewEpisodes/src/const"
	"github.com/belousandrey/NewEpisodes/src/types"
	"github.com/mmcdole/gofeed"
)

func GetNewEpisodes(resp *http.Response, lastEpisode string) (episodes []*types.Episode, last string, err error) {
	tle, err := time.Parse(constants.DateFormat, lastEpisode)
	if err != nil {
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		return
	}

	for _, e := range feed.Items {
		if e.PublishedParsed.Before(tle) {
			break
		}

		if last == "" {
			last = e.PublishedParsed.Add(time.Hour * 24).Format(constants.DateFormat)
		}

		episodes = append([]*types.Episode{types.NewEpisode(e.Title, e.Link, e.PublishedParsed.Format(constants.DateFormat))}, episodes...)
	}

	return
}
