package etool

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
)

const overwriteFlag string = "-overwrite_original"

func Write(name string, tags metadata.Tags, useSeparator, saveOriginals bool) ([]byte, error) {
	res, err := run(config.App.ToolPath(), packWriteArgs(name, tags, useSeparator, saveOriginals)...)
	if err == nil {
		log.Log("write into file", name)
	}

	return res, err
}

func packWriteArgs(name string, tags metadata.Tags, useSeparator bool, saveOriginals bool) []string {
	var args []string

	for tag, value := range tags {
		args = append(args, fmt.Sprintf("-%s=%v", tag, value))
	}

	if useSeparator {
		args = append(args, "-sep", metadata.Separator())
	}

	if !saveOriginals {
		args = append(args, overwriteFlag)
	}

	args = append(args, name)

	return args
}
