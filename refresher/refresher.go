package refresher

import (
	"sync"

	app "github.com/belousandrey/new-episodes"
	"github.com/belousandrey/new-episodes/engines/defaultengine"
	"github.com/belousandrey/new-episodes/engines/matchdaybiz"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type refresher struct {
	Podcasts     []app.Podcast
	notification bool
}

type Refresher interface {
	GetPodcasts() []app.Podcast
	SetPodcasts(podcasts []app.Podcast) Refresher
	Refresh()
	MustUpdateConfig() bool
	processPodcast(i int, podcast app.Podcast, ch chan<- app.PodcastWithEpisodes, errCh chan<- app.Podcast)
}

func NewRefresher() Refresher {
	return &refresher{}
}

func (r *refresher) SetPodcasts(podcasts []app.Podcast) Refresher {
	r.Podcasts = podcasts
	return r
}

func (r *refresher) GetPodcasts() []app.Podcast {
	return r.Podcasts
}

func (r *refresher) Refresh() {
	// wait group for all workers
	wg := new(sync.WaitGroup)

	// channel with all new episodes
	ch := make(chan app.PodcastWithEpisodes, 0)

	// read newEpisodes from ch channel
	var newEpisodes = make([]app.PodcastWithEpisodes, 0)

	// channel with errors while processing podcasts
	errCh := make(chan app.Podcast, 0)

	// read errPodcasts from errCh channel
	var problemPodcasts = make([]app.Podcast, 0)

	// go-routine to read chan with new episodes
	go func(wg *sync.WaitGroup, success <-chan app.PodcastWithEpisodes) {
		for pwe := range success {
			if pwe.LastEpisodeDate != "" {
				newEpisodes = append(newEpisodes, pwe)
				r.Podcasts[pwe.Position].Last = pwe.LastEpisodeDate
			}

			wg.Done()
		}
	}(wg, ch)

	// go-routine to read chan with problem podcasts
	go func(problems <-chan app.Podcast) {
		for pr := range problems {
			problemPodcasts = append(problemPodcasts, pr)
		}
	}(errCh)

	// start go routine for each podcast for parallel processing
	for i, e := range r.Podcasts {
		wg.Add(1)
		go r.processPodcast(i, e, ch, errCh)
	}

	wg.Wait()
	close(ch)

	// send email in case of new episodes and update config
	if len(newEpisodes) > 0 || len(problemPodcasts) > 0 {
		err := SendEmail(viper.GetString("email.to"), viper.GetStringMapString("email.from"), app.NewEmailContent(newEpisodes, problemPodcasts))
		if err != nil {
			panic(errors.Wrap(err, "send email"))
		}
	}

	if len(newEpisodes) > 0 {
		r.notification = true
	}
}

func (r *refresher) MustUpdateConfig() bool {
	return r.notification
}

func (r *refresher) processPodcast(i int, podcast app.Podcast, ch chan<- app.PodcastWithEpisodes, errCh chan<- app.Podcast) {
	resp, err := DownloadFile(podcast.Link)
	if err != nil {
		return
	}
	defer resp.Close()

	var (
		engineEpisodes []app.Episode
		last           string
	)
	switch podcast.Engine {
	case "golangshow", "changelog", "rucast", "podfm", "podster":
		engineEpisodes, last, err = defaultengine.NewEngine(podcast.Last).GetNewEpisodes(resp)
	case "matchdaybiz":
		engineEpisodes, last, err = matchdaybiz.NewEngine(podcast.Last).GetNewEpisodes(resp)
	}
	if err != nil {
		errCh <- podcast
	}

	pwe := app.NewPodcastWithEpisodes(podcast, i, last)
	pwe.Episodes = append(pwe.Episodes, engineEpisodes...)
	ch <- *pwe

	return
}
