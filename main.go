package main

import (
	"github.com/vvval/go-metadata-scanner/cmd"
	"github.com/wolfy-j/goffli/utils"
)

func main() {
	utils.Verbose = true
	cmd.Execute()
}
