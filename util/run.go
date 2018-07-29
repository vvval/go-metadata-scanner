package util

import (
	"bytes"
	"github.com/wolfy-j/goffli/utils"
	"os/exec"
)

func Run(cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	utils.Log(cmd, args...)

	errBuffer := new(bytes.Buffer)
	c.Stderr = errBuffer

	res, err := c.Output()
	if err != nil {
		utils.Printf("<red>%s</reset>", errBuffer.String())

		return []byte{}, err
	}

	return res, nil
}
