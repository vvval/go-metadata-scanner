package etool

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
)

func Write(name string, payload metadata.Payload, saveOriginals bool) ([]byte, error) {
	var args []string

	for tag, value := range payload.Tags() {
		args = append(args, fmt.Sprintf("-%s=%v", tag, value))
	}

	if payload.UseSeparator() {
		args = append(args, "-sep", metadata.Separator())
	}

	if !saveOriginals {
		args = append(args, "-overwrite_original")
	}

	args = append(args, name)

	return run(configuration.App.ToolPath(), args...)
}
