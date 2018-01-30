package types

type Podcast struct {
	Last   string
	Link   string
	Title  string
	Engine string
}

func NewPodcast(e, h, l, t string) *Podcast {
	return &Podcast{
		Title:  t,
		Link:   h,
		Engine: e,
		Last:   l,
	}
}

type Episode struct {
	Link  string
	Title string
	Date  string
}

func NewEpisode(title, link, date string) *Episode {
	return &Episode{
		Title: title,
		Link:  link,
		Date:  date,
	}
}

type PodcastWithEpisodes struct {
	Position int
	Podcast
	Episodes        []Episode
	LastEpisodeDate string
}

func NewPodcastWithEpisodes(podcast Podcast, pos int, led string) *PodcastWithEpisodes {
	return &PodcastWithEpisodes{
		pos,
		podcast,
		make([]Episode, 0),
		led,
	}
}
