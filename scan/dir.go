package scan

import (
	"github.com/vvval/go-metadata-scanner/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var files = map[string]string{}

func MustDir(directory string, extensions []string) []string {
	dirs, err := scanDir(directory, extensions)
	if err != nil {
		log.Fatalln(err)
	}

	return dirs
}

func scanDir(directory string, extensions []string) ([]string, error) {
	err := filepath.Walk(directory, visit)
	if err != nil {
		return []string{}, err
	}

	var scanned []string
	for file, extension := range files {
		if extensionMatch(extension, extensions) {
			scanned = append(scanned, file)
		}
	}

	return scanned, nil
}

func visit(path string, f os.FileInfo, _ error) error {
	if !f.IsDir() {
		path = strings.ToLower(path)
		files[path] = strings.Trim(filepath.Ext(path), ".")
	}

	return nil
}

func extensionMatch(extension string, extensions []string) bool {
	for _, ext := range extensions {
		if util.Equals(extension, ext) {
			return true
		}
	}

	return false
}
