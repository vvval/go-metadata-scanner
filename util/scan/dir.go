package scan

import (
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var files = map[string]string{}

func MustDir(directory string, extensions []string) vars.Chunk {
	files = make(map[string]string)
	dirs, err := scanDir(directory, extensions)
	if err != nil {
		log.Fatalln(err)
	}

	return dirs
}

func scanDir(directory string, extensions []string) (vars.Chunk, error) {
	err := filepath.Walk(directory, visit)
	if err != nil {
		return []string{}, err
	}

	var scanned vars.Chunk
	for file, extension := range files {
		if extensionMatch(extension, extensions) {
			scanned = append(scanned, file)
		}
	}

	return scanned, nil
}

func visit(path string, f os.FileInfo, _ error) error {
	if f == nil {
		return nil
	}

	if !f.IsDir() {
		path = strings.ToLower(path)
		files[path] = util.Extension(path)
	}

	return nil
}

func extensionMatch(extension string, extensions []string) bool {
	if len(extensions) == 0 {
		return true
	}

	for _, ext := range extensions {
		if strings.EqualFold(extension, ext) {
			return true
		}
	}

	return false
}
