package util

import (
	"bytes"
	"github.com/wolfy-j/goffli/utils"
	"os/exec"
)

func Run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	utils.Log(cmd, args...)

	errb := new(bytes.Buffer)
	c.Stderr = errb

	res, err := c.Output()
	if err != nil {
		utils.Printf("<red>%s</reset>", errb.String())
		return "", err
	}

	return string(res), nil
}
