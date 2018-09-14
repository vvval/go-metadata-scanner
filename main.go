package main

import (
	"github.com/vvval/go-metadata-scanner/cmd"
	"github.com/vvval/go-metadata-scanner/configurator"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/wolfy-j/goffli/utils"
)

func init() {
	utils.Verbose = true
	log.Visibility.Command = true
	log.Visibility.Log = true
	log.Visibility.Failure = true
	log.Visibility.Debug = true

	configurator.LoadDict()
}

func main() {
	cmd.Execute()
}
