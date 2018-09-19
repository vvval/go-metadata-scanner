package etool

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/vars"
)

var readFlags = []string{"-j", "-G", "-b"}

func Read(names vars.Chunk, fields []string) ([]byte, error) {
	return run(config.App.ToolPath(), packReadArgs(names, fields)...)
}

func packReadArgs(names vars.Chunk, fields []string) []string {
	var args = readFlags

	for _, field := range fields {
		args = append(args, fmt.Sprintf("-%s:all", field))
	}

	args = append(args, names...)

	return args
}
