package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/writecmd"
	"github.com/vvval/go-metadata-scanner/cmd/writecmd/operations"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/etool"
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
		Run: writeHandler,
	}

	rootCmd.AddCommand(cmd)
	writeFlags.Fill(cmd)
}

func writeHandler(cmd *cobra.Command, args []string) {
	files = scan.MustDir(writeFlags.Directory(), config.App.Extensions())

	if writeFlags.Append() {
		log.Log("Scan files", "\"Append\" flag is enabled")

		var poolSize, chunkSize = util.AdjustSizes(len(files), PoolSize, MinChunkSize)
		filesData = operations.ScanFiles(files.Split(chunkSize), poolSize)

		log.Log("Scanned", "\n")
	}

	var wg sync.WaitGroup
	jobs := make(chan *writecmd.Job)
	writecmd.CreatePool(&wg, poolSize, jobs, poolWorker, writeFlags.Append(), writeFlags.Originals())

	file := util.MustOpenReadonlyFile(writeFlags.Filename())
	defer file.Close()

	operations.ReadCSV(file, writeFlags.Separator(), func(filename string, payload metadata.Payload) {
		wg.Add(1)
		jobs <- writecmd.NewJob(filename, payload)
	})

	wg.Wait()
	close(jobs)

	log.Log("Writing", "done")
}

func poolWorker(job *writecmd.Job, append, originals bool) ([]byte, error) {
	filename, found := scan.Candidates(job.Filename(), files, config.App.Extensions())
	if !found {
		return []byte{}, writecmd.NoFileErr
	}

	if append {
		if file, found := findScanned(filename, &filesData); found {
			job.MergePayload(file.Tags())
		}
	}

	if !job.HasPayload() {
		return []byte{}, writecmd.SkipFileErr
	}

	payload := job.Payload()
	result, err := etool.Write(
		filename,
		payload.Tags(),
		payload.UseSeparator(),
		originals,
	)

	return result, err
}

func findScanned(filename string, files *[]vars.File) (vars.File, bool) {
	for _, file := range *files {
		if util.PathsEqual(file.Filename(), filename) {
			return file, true
		}
	}

	return vars.File{}, false
}
