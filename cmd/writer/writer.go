package writer

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util"
)

func WriteFile(name string, payload metadata.Payload, saveOriginals bool) ([]byte, error) {
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

	out, err := util.Run(config.Get().ToolPath(), args...)

	return out, err
}
