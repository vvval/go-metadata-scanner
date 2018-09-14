package log

import (
	"github.com/spiral/roadrunner/cmd/rr/utils"
	"strings"
)

var Visibility struct {
	Command, Log, Failure, Debug bool
}

var format struct {
	log, err string
}

func init() {
	format.log = "<cyan+hb>►</reset> <yellow+hb>%s</reset> <green+hb>%s</reset>\n"
	format.err = "<red+hb>%s: %s</reset>\n"
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
	utils.Printf(format.log, name, strings.Join(args, " "))
}

func logError(name string, args ...string) {
	utils.Printf(format.err, name, strings.Join(args, " "))
}
