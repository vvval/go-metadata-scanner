package scan

import (
	"github.com/vvval/go-metadata-scanner/vars"
	"path/filepath"
	"regexp"
	"strings"
)

func Candidates(filename string, files vars.Chunk, extensions []string) (string, bool) {
	endings := extEndings(filename, extensions)
	var candidates = make(map[string]bool)
	var reg = &regexp.Regexp{}

	for _, file := range files {
		for _, ending := range endings {
			if strings.EqualFold(file, ending) {
				return file, true
			}

			reg = regexp.MustCompile("^(([a-zA-Z]{1,}_)?0*)?" + regexp.QuoteMeta(filepath.Base(ending)) + "$")
			if reg.MatchString(filepath.Base(file)) {
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

	return "", false
}

func extEndings(filename string, extensions []string) []string {
	ends := []string{filename}
	for _, extension := range extensions {
		ends = append(ends, filename+"."+extension)
	}

	return ends
}
