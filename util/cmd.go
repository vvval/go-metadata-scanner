package util

import (
	"bytes"
	"errors"
	"github.com/vvval/go-metadata-scanner/log"
	"os/exec"
)

func RunCommand(cmd string, args ...string) ([]byte, error) {
	command := exec.Command(cmd, args...)

	errBuffer := new(bytes.Buffer)
	command.Stderr = errBuffer
	res, err := command.Output()

	log.Command(cmd, args...)
	if err != nil {
		return []byte{}, errors.New(errBuffer.String())
	}

	return res, nil
}