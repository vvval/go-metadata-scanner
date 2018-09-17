package writers

import (
	"github.com/vvval/go-metadata-scanner/vars"
	"os"
)

type Writer interface {
	Open(filename string, headers []string) error
	Close() error
	Write(file *vars.File) error
}

type BaseWriter struct {
	file     *os.File
	filename string
	headers  []string
}

func openFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (w *BaseWriter) closeFile() {
	if w.file != nil {
		w.file.Close()
	}
}

func NewWriter(filename string, headers []string) BaseWriter {
	return BaseWriter{
		filename: filename,
		headers:  headers,
	}
}
