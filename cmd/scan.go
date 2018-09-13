package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/etool"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/scan"
	"github.com/vvval/go-metadata-scanner/vars"
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

	for file := range scannedFiles {
		writeToFile(file)
	}
}

func writeToFile(file vars.File) {
	//fmt.Printf("+++FILE METADATA:%+v\n\n", file)
}
