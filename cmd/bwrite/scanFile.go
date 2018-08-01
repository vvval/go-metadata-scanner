package bwrite

import (
	"encoding/csv"
	"github.com/wolfy-j/goffli/utils"
	"io"
	"sync"
)

//not working
func ScanFile(reader *csv.Reader, wg sync.WaitGroup, jobs chan<- Job, in input) {
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
			utils.Printf("<red>Error on reading line %d: %s</reset>\n", i, err)
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

		wg.Add(1)
		job := Job{
			MapLineToColumns(columns, line),
			in.Directory(),
			line[0],
			in.Originals(),
			in.Append(),
		}

		jobs <- job
	}
}

func skipLine(line []string) bool {
	return len(line) == 0 || len(line[0]) == 0
}
