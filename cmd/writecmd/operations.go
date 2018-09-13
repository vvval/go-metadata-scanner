package writecmd

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/dict"
	"github.com/vvval/go-metadata-scanner/log"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"io"
	"os"
)

func WriteFile(name string, payload metadata.Payload, saveOriginals bool) ([]byte, error) {
	var args []string

	for tag, value := range payload.Tags() {
		args = append(args, fmt.Sprintf("-%s=%v", tag, value))
	}

	if payload.UseSeparator() {
		args = append(args, "-sep", metadata.Separator())
	}

	if !saveOriginals {
		args = append(args, "-overwrite_original")
	}

	args = append(args, name)

	return util.RunCommand(config.Get().ToolPath(), args...)
}

func ReadFile(file *os.File, sep rune, callback func(filename string, payload metadata.Payload)) {
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

	log.Log("Columns are:", cols...)
}
