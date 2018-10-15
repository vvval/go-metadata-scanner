package log

import (
	"strings"
)

var Visibility struct {
	Command, Log, Failure, Debug bool
}

func Debug(name string, args ...string) {
	if !Visibility.Debug {
		return
	}

	log(name, args...)
}

func Log(name string, args ...string) {
	if !Visibility.Log {
		return
	}

	log(name, args...)
}

func Command(name string, args ...string) {
	if !Visibility.Command {
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
	printf("<cyan+hb>►</reset> <yellow+hb>%s</reset> <green+hb>%s</reset>\n", name, strings.Join(args, " "))
}

func logError(name string, args ...string) {
	printf("<red+hb>%s: %s</reset>\n", name, strings.Join(args, " "))
}
