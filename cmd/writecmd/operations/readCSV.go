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

		row, err := readIntoMap(reader)
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

		callback(row[0], mapPayload(columns, row, dict))
	}
}

func readIntoMap(reader *csv.Reader) (map[int]string, error) {
	row, err := reader.Read()
	if err == nil || err == io.EOF {
		m := make(map[int]string, len(row))
		for i, v := range row {
			m[i] = v
		}

		return m, err
	}

	return nil, err
}

func skipLine(line map[int]string) bool {
	return line == nil || len(line) == 0 || len(line[0]) == 0
}

func logColumns(columns map[int]vars.Tag) {
	var cols []string
	for _, col := range columns {
		cols = append(cols, col.Key())
	}

	log.Log("Parsed columns are:", cols...)
}
