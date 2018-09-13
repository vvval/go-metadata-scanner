package scancmd

import "github.com/spf13/cobra"

type Flags struct {
	directory string
	format    string
	filename  string
}

func (f Flags) Filename() string {
	return f.filename
}

func (f Flags) Directory() string {
	return f.directory
}

func (f Flags) Format() string {
	return f.format
}

func (f *Flags) Fill(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.directory, "directory", "d", "", "Directory with files to be scanned")
	cmd.MarkFlagRequired("directory")
	cmd.Flags().StringVarP(&f.filename, "output", "o", "", "Output file name (without extension, it is specified by \"format\" flag)")
	cmd.MarkFlagRequired("filename")
	cmd.Flags().StringVarP(&f.format, "format", "f", "csv", "Output file format")
}
