package util

import (
	"os"
	"path/filepath"
	"strings"
)

var files = map[string]string{}

func Scan(directory string, extensions []string) ([]string, error) {
	err := filepath.Walk(directory, visit)
	if err != nil {
		return []string{}, err
	}

	var res = []string{}
	for file, extension := range files {
		for _, ext := range extensions {
			if extension == ext {
				res = append(res, file)
			}
		}
	}

	return res, nil
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		path = strings.ToLower(path)
		files[path] = strings.Trim(filepath.Ext(path), ".")
	}

	return nil
}
