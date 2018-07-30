package bwrite

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/cmd/config"
	"github.com/vvval/go-metadata-scanner/util"
	"strings"
)

func WriteFile(names []string, tags map[string]string) ([]byte, error) {
	var args []string
	for tag, value := range tags {
		args = append(args, fmt.Sprintf("-%s=%v", tag, convertValue(value)))
	}

	for _, name := range names {
		args = append(args, name)
	}

	out, err := util.Run(config.AppConfig().ExifToolPath, args...)

	return []byte(out), err
}

func convertValue(value string) interface{} {
	if strings.EqualFold(value, "true") {
		return true
	}

	if strings.EqualFold(value, "false") {
		return false
	}

	return value
}
