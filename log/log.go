package log

import (
	"github.com/wolfy-j/goffli/utils"
	"strings"
)

var Visibility struct {
	Command, Success, Failure, Debug bool
}

func Debug(name string, args ...string) {
	if !Visibility.Debug {
		return
	}

	utils.Log(name, args...)
}

func Log(name string, args ...string) {
	utils.Log(name, args...)
}

func Command(name string, args ...string) {
	if !Visibility.Command {
		return
	}

	utils.Log(name, args...)
}

func Success(name string, args ...string) {
	if !Visibility.Success {
		return
	}

	utils.Log(name, args...)
}

func Failure(name string, args ...string) {
	if !Visibility.Failure {
		return
	}

	logError(name, args...)
}

func logError(name string, args ...string) {
	if len(name) == 0 {
		utils.Printf("<red+hb>%s</reset>\n", strings.Join(args, " "))
	} else {
		utils.Printf("<red+hb>%s: %s</reset>\n", name, strings.Join(args, " "))
	}
}
