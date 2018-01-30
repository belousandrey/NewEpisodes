package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"sync"

	"github.com/belousandrey/NewEpisodes/src/engines/defaultengine"
	"github.com/belousandrey/NewEpisodes/src/engines/matchdaybiz"
	"github.com/belousandrey/NewEpisodes/src/types"
)

func main() {
	readConfig()

	var podcasts []types.Podcast
	err := viper.UnmarshalKey("podcasts", &podcasts)
	if err != nil {
		panic(err)
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

	if len(emailContent) > 0 {
		viper.Set("podcasts", podcasts)

		err = sendEmail(viper.GetString("email.to"), viper.GetStringMapString("email.from"), emailContent)
		if err != nil {
			panic(err)
		}

		// update config file with dates of last episodes
		if err := viper.WriteConfig(); err != nil {
			panic(fmt.Errorf("config write: %s", err))
		}
	}
}

func readConfig() {
	var config string
	pflag.StringVarP(&config, "config", "c", "../conf/conf.yaml", "path to config file")
	pflag.Parse()

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("can not read config file")
	}

	filePath := path.Join(path.Dir(filename), config)

	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
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
