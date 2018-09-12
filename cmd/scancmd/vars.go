package scancmd

import "github.com/vvval/go-metadata-scanner/metadata"

type FileData struct {
	filename string
	tags     metadata.Tags
}

type Chunk []string

var (
	Input     Flags
	Chunks    chan Chunk
	FilesData chan FileData
	PoolSize  = 10
)

type Flags struct {
	directory string
	format    string
	filename  string
}

func (input Flags) Filename() string {
	return input.filename
}

func (input Flags) Directory() string {
	return input.directory
}

func (input Flags) Format() string {
	return input.format
}
