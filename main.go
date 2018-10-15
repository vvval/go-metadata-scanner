package main

import (
	"github.com/vvval/go-metadata-scanner/cmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/util/log"
)

func init() {
	log.Visibility.Failure = true
	log.Visibility.Done = true
	log.Visibility.Debug = false   //v
	log.Visibility.Log = false     //v
	log.Visibility.Command = false //v

	config.Dict = configuration.Load(config.Dict, "./dict.yaml").(config.DictConfig)
	config.App = configuration.Load(config.App, "./app.yaml").(config.AppConfig)
	config.MSCSV = configuration.Load(config.MSCSV, "./mscsv.yaml").(config.MSCSVConfig)
}

func main() {
	cmd.Execute()
}
