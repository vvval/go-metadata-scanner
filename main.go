package main

import (
	"github.com/vvval/go-metadata-scanner/cmd"
	"github.com/vvval/go-metadata-scanner/log"
	"github.com/wolfy-j/goffli/utils"
)

func main() {
	utils.Verbose = true

	log.Visibility.Command = true
	log.Visibility.Success = true
	log.Visibility.Failure = true
	log.Visibility.Debug = true

	cmd.Execute()
}
