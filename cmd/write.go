package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/writecmd"
	"github.com/vvval/go-metadata-scanner/cmd/writecmd/operations"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/util/scan"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"sync"
)

var (
	files      vars.Chunk
	filesData  []vars.File
	writeFlags writecmd.Flags
)

const poolSize int = 20

func init() {
	var cmd = &cobra.Command{
		Use:   "write",
		Short: "Read metadata from file and writer to images",
		Long: `Read metadata from file and writer to images.
Flags file should be a CSV file with comma-separated fields (or pass custom separator via "s" flag).
First column should be reserved for file names, its name is omitted.
Other columns should be named as keywords in a dict.yaml maps section provided
for proper mapping CSV data into appropriate metadata fields`,
		Run: func(cmd *cobra.Command, args []string) {
			writeHandler(writeFlags, config.App, config.Dict, poolSize)
		},
	}

	rootCmd.AddCommand(cmd)
	writeFlags.Fill(cmd)
}

func writeHandler(flags writecmd.Flags, appConfig config.AppConfig, dictConfig config.DictConfig, poolSize int) {
	if flags.Verbosity() {
		log.Visibility.Debug = true
		log.Visibility.Log = true
		log.Visibility.Command = true
	}

	files = scan.MustDir(flags.Directory(), appConfig.Extensions())

	if flags.Append() {
		log.Log("Scanning files...", "\"Append\" flag is enabled")

		var poolSize, chunkSize = util.AdjustSizes(len(files), PoolSize, MinChunkSize)
		filesData = operations.ScanFiles(files.Split(chunkSize), poolSize)

		log.Log("Scanning completed", "\n")
	}

	var wg sync.WaitGroup
	jobs := make(chan *writecmd.Job)
	writecmd.CreatePool(&wg, poolSize, jobs, func(job *writecmd.Job) ([]byte, error) {
		return writecmd.Work(job, flags.Append(), flags.Originals(), appConfig.Extensions(), &files, &filesData)
	})

	file := util.MustOpenReadonlyFile(flags.Filename())
	defer file.Close()

	operations.ReadCSV(util.GetCSVReader(file, flags.Separator()), dictConfig, func(filename string, payload metadata.Payload) {
		wg.Add(1)
		jobs <- writecmd.NewJob(filename, payload)
	})

	wg.Wait()
	close(jobs)

	log.Done("Writing completed", "")
}
