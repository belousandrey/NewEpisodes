package main

import (
	"os"
	"path"

	app "github.com/belousandrey/new-episodes"
	"github.com/belousandrey/new-episodes/refresher"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	readConfig()

	var podcasts []app.Podcast
	err := viper.UnmarshalKey("podcasts", &podcasts)
	if err != nil {
		panic(errors.Wrap(err, "parse config file"))
	}

	r := refresher.NewRefresher().SetPodcasts(podcasts)
	r.Refresh()

	if r.MustUpdateConfig() {
		// update config file with dates of last episodes
		viper.Set("podcasts", r.GetPodcasts())

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
