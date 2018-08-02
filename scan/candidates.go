package scan

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func Candidates(filename string, files []string) {
	var candidates []string

	for _, file := range files {
		if strings.EqualFold(file, filename) {
			//the file
		}

		if strings.Index(file, filename) != -1 {
			candidates = append(candidates, file)
		}
	}

	fmt.Printf("candidates: %v = %+v\n", filename, candidates)
}

//var files map[string]string

func ScanDir(directory string, extensions []string) map[string]string {
	res, err := Dir(directory, extensions)
	if err != nil {
		log.Fatalln(err)
	}

	files := make(map[string]string)
	for _, file := range res {
		files[file] = filepath.Base(file)
	}

	return files
}

func filenameCandidates(dir, name string) []string {
	return []string{filepath.Join(dir, name)}
}
