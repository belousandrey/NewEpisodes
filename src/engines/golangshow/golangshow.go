package golangshow

import (
	"net/http"

	"github.com/belousandrey/NewEpisodes/src/engines/defaultengine"
	"github.com/belousandrey/NewEpisodes/src/types"
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
	return defaultengine.GetNewEpisodes(resp, e.lastEpisode)
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
