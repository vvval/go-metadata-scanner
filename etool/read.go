package etool

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/vars"
)

var readFlags = []string{"-j", "-G"}

func Read(names vars.Chunk) ([]byte, error) {
	var args = readFlags

	for _, field := range config.Get().Fields() {
		args = append(args, fmt.Sprintf("-%s:all", field))
	}

	args = append(args, names...)

	return run(config.Get().ToolPath(), names...)
}
