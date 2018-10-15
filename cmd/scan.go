package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd/writers"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/etool"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/util/rand"
	"github.com/vvval/go-metadata-scanner/util/scan"
	"github.com/vvval/go-metadata-scanner/vars"
	"os"
	"path/filepath"
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
	if scanFlags.Verbosity() {
		log.Visibility.Debug = true
		log.Visibility.Log = true
		log.Visibility.Command = true
	}

	log.Log("Scanning...", fmt.Sprintf("Directory is \"%s\"", util.Abs(scanFlags.Directory())))

	var files = scan.MustDir(scanFlags.Directory(), config.App.Extensions())
	poolSize, chunkSize := util.AdjustSizes(len(files), PoolSize, MinChunkSize)

	var chunks = make(chan vars.Chunk)
	var scannedFiles = make(chan vars.File)
	var wg sync.WaitGroup
	scancmd.CreatePool(
		&wg,
		poolSize,
		chunks,
		func(files vars.Chunk) ([]byte, error) {
			return etool.Read(files, config.App.Fields())
		},
		func(data []byte) {
			for _, parsed := range etool.Parse(data) {
				scannedFiles <- parsed
			}
		},
	)

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

	headers := packHeaders(config.App.Fields())
	wr, err := writers.Get(scanFlags.Format())
	if err != nil {
		logWriterFatal(err)
	}

	err = wr.Open(outputFilename, headers)
	if err != nil {
		logWriterFatal(err)
	}
	defer wr.Close()

	for file := range scannedFiles {
		file.WithRelPath(scanFlags.Directory())
		err := wr.Write(&file)
		if err != nil {
			log.Failure("CSV write", fmt.Sprintf("failed writing data for \"%s\" file", file.RelPath()))
		}
	}

	log.Done("Scanning completed", fmt.Sprintf("Output file is \"%s\" file", outputFilename))
}

func randomizeOutputFilename(path string) string {
	ext := filepath.Ext(path)
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	hash := rand.Strings(10)

	return filepath.Join(dir, base[0:len(base)-len(ext)]+"-"+hash+ext)
}

func packHeaders(fields []string) []string {
	headers := []string{"filename"}

	for _, field := range fields {
		headers = append(headers, field)
	}

	return headers
}

func logWriterFatal(err error) {
	log.Failure("Output writer", err.Error())
	os.Exit(1)
}
