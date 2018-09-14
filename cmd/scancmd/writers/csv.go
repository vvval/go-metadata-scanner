package writers

import "github.com/vvval/go-metadata-scanner/vars"

type CSVWriter struct {
	WriterProps
}

// Headers to be like: Filename, XMP, IPTC, etc...
func (w *CSVWriter) Write(files *[]vars.File) (n int, err error) {
	return 0, nil
}

func NewCSVWriter(filename string, headers []string) Writer {
	//return &Writer{filename, headers}
	return &CSVWriter{WriterProps{filename, headers}}
}
