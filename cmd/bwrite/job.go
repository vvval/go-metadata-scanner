package bwrite

import (
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
)

type Job struct {
	Line          metadata.Line
	Directory     string
	Filename      string
	SaveOriginals bool
	Append        bool
}
