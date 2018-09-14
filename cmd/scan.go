package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/etool"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/util/rand"
	"github.com/vvval/go-metadata-scanner/util/scan"
	"github.com/vvval/go-metadata-scanner/vars"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	scanFlags    scancmd.Flags
	PoolSize     = 10
	MinChunkSize = 5
)

func init() {
	// cmd represents the scan command
	var cmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan folder and write metadata into the output file.",
		Long: `Scan folder and write metadata into the output file.
By default output file is a "csv" file.`,
		Run: scanHandler,
	}

	rootCmd.AddCommand(cmd)
	scanFlags.Fill(cmd)
}

func scanHandler(cmd *cobra.Command, args []string) {
	log.Log("Scanning", util.Abs(scanFlags.Directory()))

	var files = scan.MustDir(scanFlags.Directory(), config.Get().Extensions())
	poolSize, chunkSize := util.AdjustPoolSize(PoolSize, len(files), MinChunkSize)

	var chunks = make(chan vars.Chunk)
	var scannedFiles = make(chan vars.File)
	var wg sync.WaitGroup
	scancmd.CreatePool(&wg, poolSize, chunks, etool.Read, scannedFiles)

	for _, chunk := range files.Split(chunkSize) {
		wg.Add(1)
		chunks <- chunk
	}

	go func() {
		wg.Wait()
		close(chunks)
		close(scannedFiles)
	}()

	outputFilename := randomizeOutputFilename(scanFlags.Filename())
	file, err := os.Create(outputFilename)
	if err != nil {
		log.Failure("CSV write", fmt.Sprintf("failed writing into \"%s\" file", outputFilename))
		os.Exit(1)
	}

	defer file.Close()

	headers := packHeaders()
	if scanFlags.Format() == "csv" {
		csvw := csv.NewWriter(file)
		csvw.Write(headers)
		defer csvw.Flush()

		for scannedFile := range scannedFiles {
			record := []string{scannedFile.RelPath(scanFlags.Directory())}
			for group, data := range scannedFile.PackStrings(headers) {
				for i := 1; i < len(headers); i++ {
					header := headers[i]
					if strings.EqualFold(header, group) {
						record = append(record, data)
					}
				}
			}

			err := csvw.Write(record)
			if err != nil {
				log.Failure("CSV write", fmt.Sprintf("failed writing data for \"%s\" file", scannedFile.RelPath(scanFlags.Directory())))
			}
		}
	}

	if scanFlags.Format() == "json" {
		o := map[string]map[string]interface{}{}
		for scannedFile := range scannedFiles {
			record := map[string]interface{}{}
			for k, v := range scannedFile.PackMap(headers) {
				record[k] = v
			}
			o[scannedFile.Filename()] = record
		}
		m, err := json.Marshal(o)
		if err == nil {
			file.Write(m)
		}
	}

	//for file := range scannedFiles {
	//	for k, v := range file.PackStrings(config.Get().Fields()) {
	//		fmt.Printf("group: %v\n", k)
	//		fmt.Printf("%s\n\n", v)
	//	}
	//	//writeToFile(file)
	//}

	log.Log("Scanned", fmt.Sprintf("Scan results are in the \"%s\" file", outputFilename))
}

//func getWriter(f scancmd.Flags) *writers.CSVWriter {
//	headers := []string{"filename"}
//	for _, field := range config.Get().Fields() {
//		headers = append(headers, field)
//	}
//
//	return writers.NewCSVWriter(scanFlags.Filename(), headers)
//}

//func writeToFile(file vars.File) {
//	headers := []string{scanFlags.Filename()}
//
//	for _, field := range config.Get().Fields() {
//		headers = append(headers, field)
//	}
//
//	fmt.Printf("filename %s\n format %s\nheaders %+s\n", scanFlags.Filename(), scanFlags.Format(), headers)
//}

func randomizeOutputFilename(path string) string {
	ext := filepath.Ext(path)
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	hash := rand.Strings(10)

	return filepath.Join(dir, base[0:len(base)-len(ext)]+"-"+hash+ext)
}

func packHeaders() []string {
	headers := []string{"filename"}

	for _, field := range config.Get().Fields() {
		headers = append(headers, field)
	}

	return headers
}
