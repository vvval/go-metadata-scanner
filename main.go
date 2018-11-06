package main

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/cmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"os"
	"path/filepath"
)

func init() {
	log.Visibility.Failure = true
	log.Visibility.Done = true
	log.Visibility.Debug = false   //v
	log.Visibility.Log = false     //v
	log.Visibility.Command = false //v

	dir, err := util.RootDir()
	if err != nil {
		log.Failure("System error", fmt.Sprintf("Can't detect app directory: %s", err.Error()))
		os.Exit(1)
	}

	config.Dict = configuration.Load(config.Dict, filepath.Join(dir, "dict.yaml")).(config.DictConfig)
	config.App = configuration.Load(config.App, filepath.Join(dir, "app.yaml")).(config.AppConfig)
	config.MSCSV = configuration.Load(config.MSCSV, filepath.Join(dir, "mscsv.yaml")).(config.MSCSVConfig)
}

func main() {
	cmd.Execute()
}
