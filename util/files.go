package util

import (
	"encoding/csv"
	"os"
)

func OpenReadonlyFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
}

func GetCSVReader(file *os.File, sep rune) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = sep
	reader.FieldsPerRecord = -1

	return reader
}
