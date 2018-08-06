package writeCommand

import (
	"path/filepath"
)

const defaultSeparator rune = ','

type Input struct {
	filename      string
	directory     string
	separator     string
	append        bool
	saveOriginals bool
}

func (input Input) Filename() string {
	return input.filename
}

func (input Input) Directory() string {
	if len(input.directory) != 0 {
		return input.directory
	}

	return filepath.Dir(input.filename)
}

func (input Input) Separator() rune {
	sep := []rune(input.separator)
	if len(sep) > 0 {
		return sep[0]
	}

	return defaultSeparator
}

func (input Input) Append() bool {
	return input.append
}

func (input Input) Originals() bool {
	return input.saveOriginals
}
