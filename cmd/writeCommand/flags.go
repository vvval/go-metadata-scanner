package writeCommand

import "github.com/spf13/cobra"

func FillInput(cmd *cobra.Command, i *Input) {
	cmd.Flags().StringVarP(&i.filename, "filename", "f", "", "Metadata source file name")
	cmd.MarkFlagRequired("filename")
	cmd.Flags().StringVarP(&i.separator, "sep", "s", string(defaultSeparator), "CSV file columns separator")
	cmd.Flags().StringVarP(&i.directory, "directory", "d", "", "Directory with files to be processed")
	cmd.Flags().BoolVarP(&i.append, "append", "a", false, "Append new data to existing values?")
	cmd.Flags().BoolVarP(&i.saveOriginals, "originals", "o", false, "Save original files (overwrite with new data if not set)?")
}
