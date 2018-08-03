package writer

import (
	"encoding/csv"
	"fmt"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/log"
	"io"
	"os"
)

func Read(file *os.File, sep rune, callback func(filename string, payload metadata.Payload)) {
	reader := reader(file, sep)
	var columnsFound bool
	var columns map[int]string

	var i int
	for {
		i++
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Failure("File read", fmt.Sprintf("%d row:", i), err.Error())
			continue
		}

		if !columnsFound {
			columnsFound = true
			columns = ReadColumns(row)

			var cols []string
			for _, col := range columns {
				cols = append(cols, col)
			}

			log.Log("Mapped columns are:", cols...)
			continue
		}

		if skipLine(row) {
			continue
		}

		callback(row[0], MapPayload(columns, row))
	}
}

func reader(file *os.File, sep rune) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = sep
	reader.FieldsPerRecord = -1

	return reader
}

func skipLine(line []string) bool {
	return len(line) == 0 || len(line[0]) == 0
}
