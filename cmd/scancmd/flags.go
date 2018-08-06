package scancmd

import "github.com/spf13/cobra"

func FillInput(cmd *cobra.Command, i *Flags) {
	cmd.Flags().StringVarP(&i.directory, "directory", "d", "", "Directory with files to be scanned")
	cmd.MarkFlagRequired("directory")
	cmd.Flags().StringVarP(&i.filename, "output", "o", "", "Output file name (without extension, it is specified by \"format\" flag)")
	cmd.Flags().StringVarP(&i.format, "format", "f", "csv", "Output file format")
}
