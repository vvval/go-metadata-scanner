package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
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
	var columns map[string]int

	//rowChan := make(chan int)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			fmt.Printf("csv file ended\n")
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

		//rowChan <- 1

		fmt.Printf("csv line: %T%v\n", line, line)

	}
	fmt.Printf("csv columns: %v\n", columns)
	///

	//readchan:=make(chan []byte)
	//if err := gocsv.UnmarshalToChan(file, &readchan); err != nil { // Load clients from file
	//	panic(err)
	//}

	//var m interface{}
	//gocsv.Unmarshal(file,&m)
	//fmt.Printf("%v\n",m)

	//cmdArgs := []string{}
	//
	//for _, k := range appConfig.Fields {
	//	cmdArgs = append(cmdArgs, fmt.Sprintf("-%s:all", k))
	//}
	//cmdArgs = append(cmdArgs, "-j", "-G", cmdInput.filename)
	//
	//fmt.Println("cmd args: %+v\n", cmdArgs)
	//execCmd := exec.Command(appConfig.ExifToolPath, cmdArgs...)
	//result, err := execCmd.Output()
	////fmt.Println(string(result))
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//	a := [1]map[string]interface{}{}
	//	err = json.Unmarshal([]byte(result), &a)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
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
func mapColumns(line []string) map[string]int {
	output := map[string]int{}
	for key, values := range appConfig.TagMap {
		for index, name := range line {
			// Skip empty lines and 1st column
			if index == 0 || len(name) == 0 {
				continue
			}

			// Tag map key matches
			if strings.EqualFold(name, key) {
				output[strings.Trim(key, "")] = index

				continue
			}

			// Tag map value matches
			for _, value := range values {
				if strings.EqualFold(name, value) || strings.EqualFold(name, truncateKeyPrefix(value)) {
					output[strings.Trim(key, "")] = index

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
