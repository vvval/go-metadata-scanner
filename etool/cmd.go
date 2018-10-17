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
	root, err := util.RootDir()
	if err != nil {
		return []byte{}, err
	}

	cmd := command(root, path)
	if !util.FileExists(cmd) {
		return []byte{}, errors.New("command not found " + cmd)
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

func command(root, path string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(root, path, executable+".exe")
	}

	return filepath.Join(root, path, executable)
}
