package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/viper"

	"github.com/belousandrey/NewEpisodes/src/engines/changelog"
	"github.com/belousandrey/NewEpisodes/src/engines/golangshow"
	"github.com/belousandrey/NewEpisodes/src/engines/podfm"
	"github.com/belousandrey/NewEpisodes/src/engines/rucast"
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

		if err := viper.WriteConfig(); err != nil {
			panic(fmt.Errorf("config write: %s", err))
		}

		err = sendEmail(viper.GetString("email"), emailContent)
		if err != nil {
			panic(err)
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

	switch podcast.Engine {
	case "golangshow":
		engine := golangshow.NewEngine(podcast.Last)
		res, last, err := engine.GetNewEpisodes(resp)
		if err != nil {
			return listEpisodes, newLastEpisode, err
		}
		listEpisodes = append(listEpisodes, res...)

		if last != "" {
			newLastEpisode = last
		}
	case "changelog":
		engine := changelog.NewEngine(podcast.Last)
		res, last, err := engine.GetNewEpisodes(resp)
		if err != nil {
			return listEpisodes, newLastEpisode, err
		}
		listEpisodes = append(listEpisodes, res...)

		if last != "" {
			newLastEpisode = last
		}
	case "rucast":
		engine := rucast.NewEngine(podcast.Last)
		res, last, err := engine.GetNewEpisodes(resp)
		if err != nil {
			return listEpisodes, newLastEpisode, err
		}
		listEpisodes = append(listEpisodes, res...)

		if last != "" {
			newLastEpisode = last
		}
	case "podfm":
		engine := podfm.NewEngine(podcast.Last)
		res, last, err := engine.GetNewEpisodes(resp)
		if err != nil {
			return listEpisodes, newLastEpisode, err
		}
		listEpisodes = append(listEpisodes, res...)

		if last != "" {
			newLastEpisode = last
		}
	}

	return
}
