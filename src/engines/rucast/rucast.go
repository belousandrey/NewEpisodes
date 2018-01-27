package rucast

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
