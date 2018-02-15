package main

import (
	"path"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"sync"

	"os"

	"github.com/belousandrey/NewEpisodes/src/engines/defaultengine"
	"github.com/belousandrey/NewEpisodes/src/engines/matchdaybiz"
	"github.com/belousandrey/NewEpisodes/src/types"
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

	// chan with all new episodes
	ch := make(chan types.PodcastWithEpisodes, 0)

	var emailContent = make([]types.PodcastWithEpisodes, 0)

	// go-routine to read chan with new episodes
	go func(wg *sync.WaitGroup, success <-chan types.PodcastWithEpisodes) {
		for pwe := range success {
			if pwe.LastEpisodeDate != "" {
				emailContent = append(emailContent, pwe)
				podcasts[pwe.Position].Last = pwe.LastEpisodeDate
			}

			wg.Done()
		}
	}(wg, ch)

	// go routine to process list of podcasts
	for i, e := range podcasts {
		wg.Add(1)
		go processPodcast(i, e, ch)
	}

	wg.Wait()
	close(ch)

	// send email in case of new episodes and update config
	if len(emailContent) > 0 {
		viper.Set("podcasts", podcasts)

		err = SendEmail(viper.GetString("email.to"), viper.GetStringMapString("email.from"), emailContent)
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

func processPodcast(i int, podcast types.Podcast, ch chan<- types.PodcastWithEpisodes) {
	resp, err := DownloadFile(podcast.Link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

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
		// do something
	}

	pwe := types.NewPodcastWithEpisodes(podcast, i, last)
	pwe.Episodes = append(pwe.Episodes, engineEpisodes...)
	ch <- *pwe

	return
}
