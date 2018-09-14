package operations

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/dict"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"io"
	"os"
)

func ReadCSV(file *os.File, sep rune, callback func(filename string, payload metadata.Payload)) {
	reader := util.GetCSVReader(file, sep)
	var columnsFound bool
	var columns map[int]dict.Tag

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
			columns = readColumns(row)

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

func logColumns(columns map[int]dict.Tag) {
	var cols []string
	for _, col := range columns {
		cols = append(cols, col.Key())
	}

	log.Log("Parsed columns are:", cols...)
}
