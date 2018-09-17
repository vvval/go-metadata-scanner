package main

import (
	"github.com/vvval/go-metadata-scanner/cmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/wolfy-j/goffli/utils"
)

func init() {
	utils.Verbose = true
	log.Visibility.Command = true
	log.Visibility.Log = true
	log.Visibility.Failure = true
	log.Visibility.Debug = true

	config.Dict = configuration.Load(config.Dict, "./dict.yaml").(config.DictConfig)
	config.App = configuration.Load(config.App, "./app.yaml").(config.AppConfig)
	config.MSCSV = configuration.Load(config.MSCSV, "./mscsv.yaml").(config.MSCSVConfig)
}

func main() {
	cmd.Execute()
}
