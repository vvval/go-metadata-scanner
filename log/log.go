package log

import (
	"github.com/spiral/roadrunner/cmd/rr/utils"
	"strings"
)

var Visibility struct {
	Command, Success, Failure, Debug bool
}

func Debug(name string, args ...string) {
	if !Visibility.Debug {
		return
	}

	log(name, args...)
}

func Log(name string, args ...string) {
	log(name, args...)
}

func Command(name string, args ...string) {
	if !Visibility.Command {
		return
	}

	log(name, args...)
}

func Success(name string, args ...string) {
	if !Visibility.Success {
		return
	}

	log(name, args...)
}

func Failure(name string, args ...string) {
	if !Visibility.Failure {
		return
	}

	logError(name, args...)
}

func log(name string, args ...string) {
	utils.Printf("<cyan+hb>►</reset> <yellow+hb>%s</reset> <green+hb>%s</reset>\n", name, strings.Join(args, " "))
}

func logError(name string, args ...string) {
	if len(name) == 0 {
		utils.Printf("<red+hb>%s</reset>\n", strings.Join(args, " "))
	} else {
		utils.Printf("<red+hb>%s: %s</reset>\n", name, strings.Join(args, " "))
	}
}
