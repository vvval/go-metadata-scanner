package util

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RootDir() (string, error) {
	//for debug mode: return ".",nil
	return ".", nil
	path, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(path), nil
}

func MustOpenReadonlyFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	return file
}

func GetCSVReader(file *os.File, sep rune) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = sep
	reader.FieldsPerRecord = -1

	return reader
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}

func Extension(filename string) string {
	return strings.Trim(filepath.Ext(filename), ".")
}

func Abs(path string) string {
	abs, err := filepath.Abs(path)
	if err == nil {
		return abs
	}

	return path
}

func PathsEqual(p1, p2 string) bool {
	p1 = strings.TrimRight(p1, "/\\")
	p2 = strings.TrimRight(p2, "/\\")

	if p1 == p2 {
		return true
	}

	p1 = strings.Replace(p1, "/", string(os.PathSeparator), -1)
	p1 = strings.Replace(p1, "\\", string(os.PathSeparator), -1)
	p2 = strings.Replace(p2, "/", string(os.PathSeparator), -1)
	p2 = strings.Replace(p2, "\\", string(os.PathSeparator), -1)

	return p1 == p2
}
