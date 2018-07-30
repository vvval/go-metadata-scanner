package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/bwrite"
	"github.com/vvval/go-metadata-scanner/cmd/config"
	"github.com/vvval/go-metadata-scanner/util"
	"io"
	"log"
	"os"
	"path/filepath"
)

// bwriteCmd represents the bwrite command
var bwriteCmd = &cobra.Command{
	Use:   "bwrite",
	Short: "Read metadata from file and write to images",
	Long: `Read metadata from file and write to images. Input file should be a CSV file with comma-separated fields.
First column should be reserved for file names, its name is omitted.
Other columns should be named as keywords in a config.yaml tagmap section provided
for proper mapping CSV data into appropriate metadata fields`,
	Run: bulkwrite,
}

const csvFileDefaultSeparator rune = ','

type input struct {
	filename      string
	directory     string
	separator     string
	append        bool
	saveOriginals bool
}

func (input input) sep() rune {
	sep := []rune(input.separator)
	if len(sep) > 0 {
		return sep[0]
	}

	return csvFileDefaultSeparator
}

func (input input) dir() string {
	if len(input.directory) != 0 {
		return input.directory
	}

	return filepath.Dir(input.filename)
}

var cmdInput = input{}

func init() {
	rootCmd.AddCommand(bwriteCmd)

	bwriteCmd.Flags().StringVarP(&cmdInput.filename, "filename", "f", "", "Metadata source file name")
	bwriteCmd.MarkFlagRequired("filename")
	bwriteCmd.Flags().StringVarP(&cmdInput.separator, "sep", "s", ",", "CSV file columns separator")
	bwriteCmd.Flags().StringVarP(&cmdInput.directory, "directory", "d", "", "Directory with files to be processed")
	bwriteCmd.Flags().BoolVarP(&cmdInput.append, "append", "a", false, "Append new data to existing values?")
	bwriteCmd.Flags().BoolVarP(&cmdInput.saveOriginals, "originals", "o", false, "Save original files (overwrite with new data if not set)?")
}

func bulkwrite(cmd *cobra.Command, args []string) {
	file, err := openFile(cmdInput.filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	reader, err := bwrite.Reader(file, cmdInput.sep())
	if err != nil {
		log.Fatalln(err)
	}

	var columnsLineFound bool
	var columns map[int]string

	//rowChan := make(chan int)
	//readDone, writeDone := make(chan struct{}), make(chan struct{})
	//readLine, writeLine := make(chan bool), make(chan bool)
	//
	//go func() {
	//	var countLines = 0
	//	for {
	//		select {
	//		case <-readLine:
	//			countLines++
	//		case <-writeLine:
	//			countLines --
	//			if countLines == 0 {
	//				close(done)
	//			}
	//		case <-done:
	//			return
	//		}
	//	}
	//}()

	for {
		line, err := reader.Read()
		if err == io.EOF {
			//readDone <- struct{}{}
			break
		}

		if err != nil {
			log.Fatalln(err)

			break
		}

		if !columnsLineFound {
			columnsLineFound = true
			columns = bwrite.MapColumns(line)

			fmt.Printf("%v\n\n\n", columns)
			continue
		}

		if skipLine(line) {
			continue
		}

		//fmt.Printf("line %v\n", line)
		result, err := bwrite.WriteFile(filenameCandidates(line[0]), bwrite.MapLine(columns, line))
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println(string(result))
		}
	}

	files, err := util.Scan(cmdInput.dir(), config.AppConfig().Extensions)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("files %s\n", files)
}

func openFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func skipLine(line []string) bool {
	return len(line) == 0 || len(line[0]) == 0
}

func filenameCandidates(name string) []string {
	return []string{name}
}
