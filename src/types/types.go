package types

import "io"

// Episoder - interface for engines
type Episoder interface {
	GetNewEpisodes(resp io.Reader) (episodes []Episode, last string, err error)
}

// Podcast - config file representation for podcast
type Podcast struct {
	// Last - string with last episode date
	Last string
	// Link - link to RSS or HTML page with episodes list
	Link string
	// Title - human readable name
	Title string
	// Engine - way of podcast processing
	Engine string
}

// Episode - one episode of podcast
type Episode struct {
	// Link - URL to the episode file
	Link string
	// Title - name of the episode
	Title string
	// Date - string with episode date
	Date string
}

// PodcastWithEpisodes - data structure for email template
type PodcastWithEpisodes struct {
	// Position - position in config file
	Position int
	Podcast
	// Episodes - list with new episodes
	Episodes []Episode
	// LastEpisodeDate - string with last episode date
	LastEpisodeDate string
}

// NewPodcastWithEpisodes - create new data structure
func NewPodcastWithEpisodes(podcast Podcast, pos int, led string) *PodcastWithEpisodes {
	return &PodcastWithEpisodes{
		pos,
		podcast,
		make([]Episode, 0),
		led,
	}
}
