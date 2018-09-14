package main

import (
	"github.com/vvval/go-metadata-scanner/cmd"
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

	configuration.LoadAll()
}

func main() {
	cmd.Execute()

	//fmt.Printf("old %+v\nnew %+v\n\n\n", config.Get(), configurator.App)
}
