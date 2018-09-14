package writers

import "github.com/vvval/go-metadata-scanner/vars"

type JSONWriter struct {
	WriterProps
}

// Headers to be like: Filename, XMP, IPTC, etc...
func (w *JSONWriter) Write(files *[]vars.File) (n int, err error) {
	return 0, nil
}

func NewJSONWriter(filename string, headers []string) Writer {
	//return &Writer{filename, headers}
	return &JSONWriter{WriterProps{filename, headers}}
}
