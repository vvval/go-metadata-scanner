package writer

import (
	"github.com/spf13/cobra"
	"path/filepath"
)

const csvFileDefaultSeparator rune = ','

type input struct {
	filename      string
	directory     string
	separator     string
	append        bool
	saveOriginals bool
}

func (input input) Filename() string {
	return input.filename
}

func (input input) Directory() string {
	if len(input.directory) != 0 {
		return input.directory
	}

	return filepath.Dir(input.filename)
}

func (input input) Separator() rune {
	sep := []rune(input.separator)
	if len(sep) > 0 {
		return sep[0]
	}

	return csvFileDefaultSeparator
}

func (input input) Append() bool {
	return input.append
}

func (input input) Originals() bool {
	return input.saveOriginals
}

var i input

func InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.filename, "filename", "f", "", "Metadata source file name")
	cmd.MarkFlagRequired("filename")
	cmd.Flags().StringVarP(&i.separator, "sep", "s", string(csvFileDefaultSeparator), "CSV file columns separator")
	cmd.Flags().StringVarP(&i.directory, "directory", "d", "", "Directory with files to be processed")
	cmd.Flags().BoolVarP(&i.append, "append", "a", false, "Append new data to existing values?")
	cmd.Flags().BoolVarP(&i.saveOriginals, "originals", "o", false, "Save original files (overwrite with new data if not set)?")
}

func Input() input {
	return i
}
