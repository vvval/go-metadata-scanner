package writers

import "github.com/vvval/go-metadata-scanner/vars"

type CSVWriter struct {
	WriterProps
}

// Headers to be like: Filename, XMP, IPTC, etc...
func (w *CSVWriter) Write(files *[]vars.File) {

}

func NewCSVWriter(filename string, headers []string) *CSVWriter {
	//return &Writer{filename, headers}
	return &CSVWriter{WriterProps{filename, headers}}
}
