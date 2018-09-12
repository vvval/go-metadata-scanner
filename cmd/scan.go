package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/scan"
	"github.com/vvval/go-metadata-scanner/util"
	"math"
	"sync"
)

func init() {
	// cmd represents the scan command
	var cmd = &cobra.Command{
		Use:   "scan",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: scanHandler,
	}

	rootCmd.AddCommand(cmd)
	scancmd.FillInput(cmd, &scancmd.Input)

	scancmd.Chunks = make(chan scancmd.Chunk)
	scancmd.FilesData = make(chan scancmd.FileData, 10)
	scancmd.CreatePool(&wg, scancmd.PoolSize, scancmd.Chunks, ScanFiles, scancmd.FilesData)
}

func ScanFiles(names scancmd.Chunk) ([]byte, error) {
	return util.RunCommand(config.Get().ToolPath(), packArgs(names)...)
}

// Pack required flags, extract fields and names
func packArgs(names []string) []string {
	var args = []string{"-j", "-G"}

	for _, field := range config.Get().Fields() {
		args = append(args, fmt.Sprintf("-%s:all", field))
	}

	return append(args, names...)
}

func scanHandler(cmd *cobra.Command, args []string) {
	GetFiles(scan.MustDir(scancmd.Input.Directory(), config.Get().Extensions()), scancmd.PoolSize, &wg, scancmd.Chunks)

	go func() {
		wg.Wait()
		close(scancmd.Chunks)
		close(scancmd.FilesData)
	}()

	for file := range scancmd.FilesData {
		writeToFile(file)
	}
}

func writeToFile(file scancmd.FileData) {
	fmt.Printf("THIS IS THE SCANNED FILE METADATA:%+v\n\n", file)
}

func GetFiles(files []string, poolSize int, wg *sync.WaitGroup, out chan<- scancmd.Chunk) {
	var chunkSize = int(math.Ceil(float64(len(files) / poolSize)))

	for i := 0; i < len(files); i += chunkSize {
		end := i + chunkSize
		if end > len(files) {
			end = len(files)
		}

		wg.Add(1)

		chunk := files[i:end]
		if len(chunk) > 0 {
			out <- chunk
		}
	}
}
