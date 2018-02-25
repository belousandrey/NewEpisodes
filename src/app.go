package main

import (
	"path"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"sync"

	"os"

	"github.com/belousandrey/new-episodes/src/engines/defaultengine"
	"github.com/belousandrey/new-episodes/src/engines/matchdaybiz"
	"github.com/belousandrey/new-episodes/src/types"
	"github.com/pkg/errors"
)

func main() {
	readConfig()

	var podcasts []types.Podcast
	err := viper.UnmarshalKey("podcasts", &podcasts)
	if err != nil {
		panic(errors.Wrap(err, "parse config file"))
	}

	// wait group for all workers
	wg := new(sync.WaitGroup)

	// channel with all new episodes
	ch := make(chan types.PodcastWithEpisodes, 0)

	// read newEpisodes from ch channel
	var newEpisodes = make([]types.PodcastWithEpisodes, 0)

	// channel with errors while processing podcasts
	errCh := make(chan types.Podcast, 0)

	// read errPodcasts from errCh channel
	var problemPodcasts = make([]types.Podcast, 0)

	// go-routine to read chan with new episodes
	go func(wg *sync.WaitGroup, success <-chan types.PodcastWithEpisodes) {
		for pwe := range success {
			if pwe.LastEpisodeDate != "" {
				newEpisodes = append(newEpisodes, pwe)
				podcasts[pwe.Position].Last = pwe.LastEpisodeDate
			}

			wg.Done()
		}
	}(wg, ch)

	// go-routine to read chan with problem podcasts
	go func(problems <-chan types.Podcast) {
		for pr := range problems {
			problemPodcasts = append(problemPodcasts, pr)
		}
	}(errCh)

	// start go routine for each podcast for parallel processing
	for i, e := range podcasts {
		wg.Add(1)
		go processPodcast(i, e, ch, errCh)
	}

	wg.Wait()
	close(ch)

	// send email in case of new episodes and update config
	if len(newEpisodes) > 0 || len(problemPodcasts) > 0 {
		viper.Set("podcasts", podcasts)

		err = SendEmail(viper.GetString("email.to"), viper.GetStringMapString("email.from"), types.NewEmailContent(newEpisodes, problemPodcasts))
		if err != nil {
			panic(errors.Wrap(err, "send email"))
		}

		// update config file with dates of last episodes
		if err := viper.WriteConfig(); err != nil {
			panic(errors.Wrap(err, "config write"))
		}
	}
}

func readConfig() {
	var config string
	pflag.StringVarP(&config, "config", "c", "../conf/conf.yaml", "path to config file")
	pflag.Parse()

	var filePath string
	if path.IsAbs(config) {
		filePath = config
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			panic(errors.Wrap(err, "get working directory"))
		}

		filePath = path.Join(pwd, config)
	}

	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "read config file"))
	}
}

func processPodcast(i int, podcast types.Podcast, ch chan<- types.PodcastWithEpisodes, errCh chan<- types.Podcast) {
	resp, err := DownloadFile(podcast.Link)
	if err != nil {
		return
	}
	defer resp.Close()

	var (
		engineEpisodes []types.Episode
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

	pwe := types.NewPodcastWithEpisodes(podcast, i, last)
	pwe.Episodes = append(pwe.Episodes, engineEpisodes...)
	ch <- *pwe

	return
}
