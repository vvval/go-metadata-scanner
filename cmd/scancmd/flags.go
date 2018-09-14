package scancmd

import (
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/util"
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
	return filepath.ToSlash(f.directory)
}

func (f Flags) Filename() string {
	filename := f.filename[0:len(f.filename)-len(f.ext())] + "." + f.Format()

	return filepath.Join(f.Directory(), filename)
}

func (f Flags) Format() string {
	ext := strings.ToLower(util.Extension(f.filename))
	if len(ext) == 0 {
		ext = defaultFormat
	}

	return ext
}

func (f Flags) ext() string {
	return filepath.Ext(f.filename)
}

func (f *Flags) Fill(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.directory, "directory", "d", "", "Directory with files to be scanned")
	cmd.MarkFlagRequired("directory")
	cmd.Flags().StringVarP(&f.filename, "filename", "f", "", "Output file (with extension)")
	cmd.MarkFlagRequired("filename")
}
