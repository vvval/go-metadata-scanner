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

func init() {
	var bwriteCmd = &cobra.Command{
		Use:   "bwrite",
		Short: "Read metadata from file and write to images",
		Long: `Read metadata from file and write to images. Input file should be a CSV file with comma-separated fields.
First column should be reserved for file names, its name is omitted.
Other columns should be named as keywords in a config.yaml tagmap section provided
for proper mapping CSV data into appropriate metadata fields`,
		Run: bulkwrite,
	}

	rootCmd.AddCommand(bwriteCmd)
	bwrite.InitFlags(bwriteCmd)
}

func bulkwrite(cmd *cobra.Command, args []string) {
	input := bwrite.Input()
	file, err := openFile(input.Filename())
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	reader, err := bwrite.Reader(file, input.Separator())
	if err != nil {
		log.Fatalln(err)
	}

	var columnsLineFound bool
	var columns map[int]string

	//var lines map[string]interface{}

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

			fmt.Printf("%v\n", columns)
			continue
		}

		if skipLine(line) {
			continue
		}

		fmt.Printf("line %v\n", line)
		result, err := bwrite.WriteFile(
			filenameCandidates(input.Directory(), line[0]),
			bwrite.MapLineToColumns(columns, line),
			input.Originals(),
			input.Append(),
		)

		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println(string(result))
		}
	}

	files, err := util.Scan(input.Directory(), config.AppConfig().Extensions)
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

func filenameCandidates(dir, name string) []string {
	return []string{filepath.Join(dir, name)}
}
