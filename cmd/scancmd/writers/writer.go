package writers

import "github.com/vvval/go-metadata-scanner/vars"

type Writer interface {
	Write(files *[]vars.File) (n int, err error)
}

type WriterProps struct {
	filename string
	headers  []string
}

// Headers to be like: Filename, XMP, IPTC, etc...
//func (w *Writer) Write(files *[]vars.File) {
//
//}
