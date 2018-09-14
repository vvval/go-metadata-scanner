package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd/writers"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/etool"
	"github.com/vvval/go-metadata-scanner/util"
	log2 "github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/util/scan"
	"github.com/vvval/go-metadata-scanner/vars"
	"log"
	"os"
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

	//writer := getWriter(scanFlags)
	file, err := os.Create(scanFlags.Filename())
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	headers := packHeaders()
	w := csv.NewWriter(file)
	w.Write(headers)
	defer w.Flush()

	for file := range scannedFiles {
		record := []string{file.RelPath(scanFlags.Directory())}
		for group, data := range file.Pack(headers) {
			for i := 1; i < len(headers); i++ {
				header := headers[i]
				//fmt.Printf("headers: %+v\n%d [%s]\n", headers, i, header)
				if strings.EqualFold(header, group) {
					record = append(record, data)
				}
			}
			//for i, header := range headers {
			//	if strings.EqualFold(header, group) {
			//		record[i] = data
			//	}
			//}
		}

		err := w.Write(record)
		if err != nil {
			log.Fatalln(err)
		}
	}

	for file := range scannedFiles {
		//_=file.Pack(config.Get().Fields())
		for k, v := range file.Pack(config.Get().Fields()) {
			fmt.Printf("group: %v\n", k)
			fmt.Printf("%s\n\n", v)
		}
		//writeToFile(file)
	}

	log2.Success("Scanned", fmt.Sprintf("Scan results are in the \"%s\" file", scanFlags.Filename()))
}

func getWriter(f scancmd.Flags) *writers.CSVWriter {
	headers := []string{scanFlags.Filename()}
	for _, field := range config.Get().Fields() {
		headers = append(headers, field)
	}

	return writers.NewCSVWriter(scanFlags.Filename(), headers)
}

func writeToFile(file vars.File) {
	headers := []string{scanFlags.Filename()}

	for _, field := range config.Get().Fields() {
		headers = append(headers, field)
	}

	fmt.Printf("filename %s\n format %s\nheaders %+s\n", scanFlags.Filename(), scanFlags.Format(), headers)
}

func packHeaders() []string {
	headers := []string{"filename"}

	for _, field := range config.Get().Fields() {
		headers = append(headers, field)
	}

	return headers
}
