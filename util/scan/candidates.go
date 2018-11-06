package scan

import (
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/vars"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func Candidates(filename string, files *vars.Chunk, extensions []string) (string, bool) {
	endings := extEndings(filename, extensions)
	var candidates = make(map[string]bool)

	for _, file := range *files {
		for _, ending := range endings {
			if strings.EqualFold(file, ending) {
				return file, true
			}

			if filesMatch(filepath.Base(ending), filepath.Base(file)) {
				candidates[file] = true
			}
		}
	}

	var values []string
	for k := range candidates {
		values = append(values, k)
	}

	if len(values) == 1 {
		return values[0], true
	}

	if len(values) == 0 {
		log.Debug("No candidates for " + filename)
	} else {
		log.Debug("Too many candidates for "+filename, values...)
	}

	return "", false
}

func filesMatch(ending, file string) bool {
	if runtime.GOOS == "windows" {
		ending = strings.ToLower(ending)
		file = strings.ToLower(file)
	}

	var reg = &regexp.Regexp{}
	reg = regexp.MustCompile("^(([a-zA-Z]{1,}_)?0*)?" + regexp.QuoteMeta(ending) + "$")

	return reg.MatchString(file)
}

func extEndings(filename string, extensions []string) []string {
	ends := []string{filename}
	for _, extension := range extensions {
		ends = append(ends, filename+"."+extension)
	}

	return ends
}
