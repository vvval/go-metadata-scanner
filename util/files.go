package util

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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

	p1to := filepath.ToSlash(p1)
	p1from := filepath.FromSlash(p1)

	p2to := filepath.ToSlash(p2)
	p2from := filepath.FromSlash(p2)

	if p1to == p2to || p1to == p2from || p1from == p2to || p1from == p2from {
		return true
	}

	return false
}
