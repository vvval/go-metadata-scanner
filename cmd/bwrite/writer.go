package bwrite

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/cmd/config"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/util"
)

func WriteFile(names []string, line metadata.Line, saveOriginals bool) ([]byte, error) {
	var args []string

	for tag, value := range line.Tags() {
		args = append(args, fmt.Sprintf("-%s=%v", tag, value))
	}

	if line.UseSeparator() {
		args = append(args, fmt.Sprintf("-sep %s", metadata.Separator()))
	}

	if !saveOriginals {
		args = append(args, "-overwrite_original")
	}

	for _, name := range names {
		args = append(args, name)
	}

	out, err := util.Run(config.AppConfig().ExifToolPath, args...)

	return out, err
}
