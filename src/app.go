package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/viper"

	"github.com/belousandrey/NewEpisodes/src/engines/defaultengine"
	"github.com/belousandrey/NewEpisodes/src/engines/matchdaybiz"
	"github.com/belousandrey/NewEpisodes/src/types"
)

const (
	configFile = "../conf/conf.yaml"
)

func main() {
	readConfig()

	var podcasts []types.Podcast
	err := viper.UnmarshalKey("podcasts", &podcasts)
	if err != nil {
		panic(err)
	}

	var emailContent = make([]*types.PodcastWithEpisodes, 0)
	for i, e := range podcasts {
		episodes, nle, err := processPodcast(e)
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err.Error())
		}

		if len(episodes) > 0 {
			pwe := types.NewPodcastWithEpisodes(e)

			for _, e := range episodes {
				pwe.Episodes = append(pwe.Episodes, *e)
			}

			emailContent = append(emailContent, pwe)
		}

		if nle != "" {
			podcasts[i].Last = nle
		}
	}

	if len(podcasts) > 0 {
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
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("can not read config file")
	}

	filePath := path.Join(path.Dir(filename), configFile)

	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func processPodcast(podcast types.Podcast) (listEpisodes []*types.Episode, newLastEpisode string, err error) {
	resp, err := DownloadFile(podcast.Link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var (
		res  []*types.Episode
		last string
	)
	switch podcast.Engine {
	case "golangshow", "changelog", "rucast", "podfm", "podster":
		res, last, err = defaultengine.NewEngine(podcast.Last).GetNewEpisodes(resp)
	case "matchdaybiz":
		res, last, err = matchdaybiz.NewEngine(podcast.Last).GetNewEpisodes(resp)
	}
	if err != nil {
		return listEpisodes, newLastEpisode, err
	}
	listEpisodes = append(listEpisodes, res...)

	if last != "" {
		newLastEpisode = last
	}

	return
}
