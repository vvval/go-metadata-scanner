package etool

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/vars"
)

var readFlags = []string{"-j", "-G", "-b"}

func Read(names vars.Chunk, fields []string) ([]byte, error) {
	res, err := run(config.App.ToolPath(), packReadArgs(names, fields)...)
	if err == nil {
		log.Log("Read from file(s)", names...)
	}

	return res, err
}

func packReadArgs(names vars.Chunk, fields []string) []string {
	var args = readFlags

	for _, field := range fields {
		args = append(args, fmt.Sprintf("-%s:all", field))
	}

	args = append(args, names...)

	return args
}
