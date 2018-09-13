package writecmd

import (
	"github.com/spf13/cobra"
	"path/filepath"
)

const defaultSeparator rune = ','

type Flags struct {
	filename      string
	directory     string
	separator     string
	append        bool
	saveOriginals bool
}

func (f Flags) Filename() string {
	return f.filename
}

func (f Flags) Directory() string {
	if len(f.directory) != 0 {
		return f.directory
	}

	return filepath.Dir(f.filename)
}

func (f Flags) Separator() rune {
	sep := []rune(f.separator)
	if len(sep) > 0 {
		return sep[0]
	}

	return defaultSeparator
}

func (f Flags) Append() bool {
	return f.append
}

func (f Flags) Originals() bool {
	return f.saveOriginals
}

func (f *Flags) Fill(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.filename, "filename", "f", "", "Metadata source file name")
	cmd.MarkFlagRequired("filename")
	cmd.Flags().StringVarP(&f.separator, "sep", "s", string(defaultSeparator), "CSV file columns separator")
	cmd.Flags().StringVarP(&f.directory, "directory", "d", "", "Directory with files to be processed")
	cmd.Flags().BoolVarP(&f.append, "append", "a", false, "Append new data to existing values?")
	cmd.Flags().BoolVarP(&f.saveOriginals, "originals", "o", false, "Save original files (overwrite with new data if not set)?")
}
