package writers

import (
	"encoding/csv"
	"encoding/json"
	"github.com/vvval/go-metadata-scanner/vars"
	"strings"
)

type CSVWriter struct {
	BaseWriter
	csv *csv.Writer
}

// Headers to be like: Filename, XMP, IPTC, etc...
func (w *CSVWriter) Write(file *vars.File) error {
	record := []string{file.RelPath()}
	for group, data := range packStrings(file, w.headers) {
		for i := 1; i < len(w.headers); i++ {
			header := w.headers[i]
			if strings.EqualFold(header, group) {
				record = append(record, data)
			}
		}
	}

	return w.csv.Write(record)
}

func (w *CSVWriter) Open(filename string, headers []string) error {
	w.BaseWriter = NewWriter(filename, headers)

	file, err := openFile(w.filename)
	if err != nil {
		return err
	}

	w.file = file
	w.csv = csv.NewWriter(file)
	w.csv.Write(headers)

	return nil
}

func (w *CSVWriter) Close() error {
	if w.csv != nil {
		w.csv.Flush()
	}

	closeFile(w.file)

	return nil
}

func packStrings(file *vars.File, headers []string) map[string]string {
	output := make(map[string]string)

	for header, value := range tagsByGroups(file, headers) {
		packed, err := json.Marshal(value)
		if err == nil {
			output[header] = string(packed)
		}
	}

	return output
}
