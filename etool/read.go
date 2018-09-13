package etool

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
)

var readFlags = []string{"-j", "-G"}

func Read(names vars.Chunk) ([]byte, error) {
	return util.RunCommand(config.Get().ToolPath(), packArgs(names)...)
}

// Pack required flags, extract fields and names
func packArgs(names []string) []string {
	var args = readFlags

	for _, field := range config.Get().Fields() {
		args = append(args, fmt.Sprintf("-%s:all", field))
	}

	return append(args, names...)
}
