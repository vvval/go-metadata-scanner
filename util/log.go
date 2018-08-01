package util

import (
	"github.com/wolfy-j/goffli/utils"
	"strings"
)

func LogError(name string, args ...string) {
	if len(name) == 0 {
		utils.Printf("<red+hb>%s</reset>\n", strings.Join(args, " "))
	} else {
		utils.Printf("<red+hb>%s %s</reset>\n", name, strings.Join(args, " "))
	}
}
