package writers

import (
	"encoding/json"
	"github.com/vvval/go-metadata-scanner/vars"
)

type JSONWriter struct {
	BaseWriter
	buf map[string]map[string]interface{}
}

// Headers to be like: Filename, XMP, IPTC, etc...
func (w *JSONWriter) Write(file *vars.File) error {
	w.buf[file.RelPath()] = packJSONLine(file, w.headers)

	return nil
}

func (w *JSONWriter) Open(filename string, headers []string) error {
	w.BaseWriter = NewWriter(filename, headers)

	file, err := openFile(w.filename)
	if err != nil {
		return err
	}

	w.buf = make(map[string]map[string]interface{})
	w.file = file

	return nil
}

func (w *JSONWriter) Close() error {
	packed, err := json.MarshalIndent(w.buf, "", "  ")
	if err != nil {
		return err
	}

	w.file.Write(packed)
	closeFile(w.file)

	return nil
}

func packJSONLine(file *vars.File, headers []string) map[string]interface{} {
	record := map[string]interface{}{}
	for k, v := range tagsByGroups(file, headers) {
		record[k] = v
	}

	return record
}
