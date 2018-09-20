package operations

import (
	"encoding/csv"
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"io"
)

func ReadCSV(reader *csv.Reader, dict config.DictConfig, callback func(filename string, payload metadata.Payload)) {
	var columnsFound bool
	var columns map[int]vars.Tag

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
			columns = readColumns(row, dict)

			logColumns(columns)
			continue
		}

		if skipLine(row) {
			continue
		}

		callback(row[0], mapPayload(columns, row))
	}
}

func skipLine(line []string) bool {
	return len(line) == 0 || len(line[0]) == 0
}

func logColumns(columns map[int]vars.Tag) {
	var cols []string
	for _, col := range columns {
		cols = append(cols, col.Key())
	}

	log.Log("Parsed columns are:", cols...)
}
