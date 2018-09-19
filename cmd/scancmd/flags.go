package scancmd

import (
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

const defaultFormat string = "csv"

type Flags struct {
	directory string
	format    string
	filename  string
}

func (f Flags) Directory() string {
	return f.directory
}

func (f Flags) Filename() string {
	return filepath.Join(f.Directory(), f.filename+"."+f.ext())
}

func (f Flags) Format() string {
	return f.format
}

func (f Flags) ext() string {
	if len(f.format) == 0 || strings.EqualFold(f.format, "mscsv") {
		return defaultFormat
	}

	return f.format
}

func (f *Flags) Fill(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.directory, "directory", "d", "", "Directory with files to be scanned")
	cmd.MarkFlagRequired("directory")
	cmd.Flags().StringVarP(&f.filename, "output", "o", "", "Output file (without extension)")
	cmd.MarkFlagRequired("output")
	cmd.Flags().StringVarP(&f.format, "format", "f", "csv", "Output file format")
}
