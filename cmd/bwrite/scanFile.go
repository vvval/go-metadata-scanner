package bwrite

import (
	"encoding/csv"
	"fmt"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/wolfy-j/goffli/utils"
	"io"
)

func ScanFile(reader *csv.Reader, callback func(filename string, line metadata.Line)) {
	var columnsLineFound bool
	var columns map[int]string

	var i int
	for {
		i++
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			util.Error("File read", fmt.Sprintf("%d line:", i), err.Error())
			continue
		}

		if !columnsLineFound {
			columnsLineFound = true
			columns = MapColumns(line)

			var cols []string
			for _, col := range columns {
				cols = append(cols, col)
			}

			utils.Log("Mapped columns are:", cols...)
			continue
		}

		if skipLine(line) {
			continue
		}

		callback(line[0], MapLineToColumns(columns, line))
	}
}

func skipLine(line []string) bool {
	return len(line) == 0 || len(line[0]) == 0
}
