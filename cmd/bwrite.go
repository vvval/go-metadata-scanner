package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/util"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

type input struct {
	filename      string
	directory     string
	append        bool
	saveOriginals bool
}

var cmdInput = input{}

func init() {
	rootCmd.AddCommand(bwriteCmd)

	bwriteCmd.Flags().StringVarP(&cmdInput.filename, "filename", "f", "", "Metadata source file name")
	bwriteCmd.MarkFlagRequired("filename")
	bwriteCmd.Flags().StringVarP(&cmdInput.directory, "directory", "d", "", "Directory with files to be processed")
	bwriteCmd.Flags().BoolVarP(&cmdInput.append, "append", "a", false, "Append new data to existing values?")
	bwriteCmd.Flags().BoolVarP(&cmdInput.saveOriginals, "originals", "o", false, "Save original files (overwrite with new data if not set)?")
}

func bulkwrite(cmd *cobra.Command, args []string) {
	prepareInput(&cmdInput)

	file, err := os.OpenFile(cmdInput.filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	reader := csvReader(file)

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
			columns = mapColumns(line)

			continue
		}

		if len(line) == 0 || len(line[0]) == 0 {
			continue
		}

		result, err := writeLine(filenameCandidates(line[0]), mapLineData(columns, line))
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println(string(result))
		}
	}

	//var cmdArgs = []string{}
	//for _, k := range appConfig.Fields {
	//	cmdArgs = append(cmdArgs, fmt.Sprintf("-%s:all", k), "-jfif:all")
	//}
	//fmt.Printf("%v",cmdArgs)
	//cmdArgs = append(cmdArgs, `-Headline=my head is "spinning"! ~=hello`, "/Users/xavier/Documents/123.jpg")
	//fmt.Printf("%v",cmdArgs)
	//execCmd := exec.Command(appConfig.ExifToolPath, cmdArgs...)
	//result, err := execCmd.Output()
	//fmt.Println(string(result))
	//fmt.Printf("%v\n",columns)
}

// Do some preparations to a user input
func prepareInput(input *input) {
	if len(input.directory) == 0 {
		input.directory = filepath.Dir(input.filename)
	}
}

func csvReader(csvFile *os.File) *csv.Reader {
	reader := csv.NewReader(csvFile)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	return reader
}

// Map columns to a known tag map
func mapColumns(line []string) map[int]string {
	output := map[int]string{}
	for key, values := range appConfig.TagMap {
		for index, name := range line {
			// Skip empty lines and 1st column
			if index == 0 || len(name) == 0 {
				continue
			}

			// Tag map key matches
			if strings.EqualFold(name, key) {
				output[index] = strings.Trim(key, "")

				continue
			}

			// Tag map value matches
			for _, value := range values {
				if strings.EqualFold(name, value) || strings.EqualFold(name, truncateKeyPrefix(value)) {
					output[index] = strings.Trim(key, "")

					break
				}
			}
		}
	}

	return output
}

// Cut <group:> prefix if found
func truncateKeyPrefix(key string) string {
	prefixEnding := strings.Index(key, ":")
	if prefixEnding == -1 {
		return key
	}

	runes := []rune(key)

	return string(runes[prefixEnding+1:])
}

func filenameCandidates(name string) []string {
	return []string{name}
}

func mapLineData(columns map[int]string, data []string) map[string]string {
	output := map[string]string{}

	for index, value := range data {
		key, ok := columns[index]
		// Unmapped key, skip
		if !ok {
			continue
		}

		// Unknown tag
		tags, ok := (appConfig.TagMap)[key]
		if !ok {
			continue
		}

		for _, tag := range tags {
			output[tag] = value
		}
	}

	return output
}

func writeLine(names []string, tags map[string]string) ([]byte, error) {
	var args []string
	for tag, value := range tags {
		args = append(args, fmt.Sprintf("-%s=%v", tag, convertValue(value)))
	}

	for _, name := range names {
		args = append(args, name)
	}

	out, err := util.Run(appConfig.ExifToolPath, args...)

	return []byte(out), err
}

func convertValue(value string) interface{} {
	if strings.EqualFold(value, "true") {
		return true
	}

	if strings.EqualFold(value, "false") {
		return false
	}

	return value
}
