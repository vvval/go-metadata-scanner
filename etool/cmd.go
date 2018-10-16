package etool

import (
	"bytes"
	"errors"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"os/exec"
	"path/filepath"
	"runtime"
)

const executable = "exiftool"

func run(path string, args ...string) ([]byte, error) {
	cmd := command(path)
	if !util.FileExists(cmd) {
		return []byte{}, errors.New("command not found")
	}

	command := exec.Command(cmd, args...)

	errBuffer := new(bytes.Buffer)
	command.Stderr = errBuffer
	res, err := command.Output()

	log.Command(executable, args...)
	if err != nil {
		return []byte{}, errors.New(errBuffer.String())
	}

	return res, nil
}

func command(path string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(path, executable+".exe")
	}

	return filepath.Join(path, executable)
}
