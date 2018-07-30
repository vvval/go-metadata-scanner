package bwrite

import (
	"encoding/csv"
	"os"
)

func Reader(file *os.File, sep rune) (*csv.Reader, error) {
	reader := csv.NewReader(file)
	reader.Comma = sep
	reader.FieldsPerRecord = -1

	return reader, nil
}
