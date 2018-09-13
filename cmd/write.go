package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/writecmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/log"
	"github.com/vvval/go-metadata-scanner/scan"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"sync"
)

type Job struct {
	filename string
	payload  metadata.Payload
}

var (
	wg           sync.WaitGroup
	jobs         chan *Job
	files        []string
	appFilesData []vars.File
	skipFileErr  = errors.New("skipFileErr")
	noFileErr    = errors.New("noFileErr")
	writeInput   writecmd.Flags
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
	writeInput.Fill(cmd)
}

func writeHandler(cmd *cobra.Command, args []string) {
	if writeInput.Append() {
		log.Log("Scan files for appending", "")

		var files = scan.MustDir(writeInput.Directory(), config.Get().Extensions())
		var poolSize, chunkSize = util.AdjustPoolSize(PoolSize, len(files), MinChunkSize)
		appFilesData = writecmd.Scan(files.Split(chunkSize), poolSize)

		log.Success("Scanned", "\n")
	}

	jobs = make(chan *Job)
	initPool(poolSize, jobs, poolWorker)

	file := util.MustOpenReadonlyFile(writeInput.Filename())
	defer file.Close()

	writecmd.ReadFile(file, writeInput.Separator(), func(filename string, payload metadata.Payload) {
		wg.Add(1)
		jobs <- &Job{filename, payload}
	})

	wg.Wait()
	close(jobs)

	log.Success("Writing", "done")
}

func initPool(poolSize int, jobs <-chan *Job, poolWorkerFunc func(job *Job, append, originals bool) ([]byte, error)) {
	// worker pool
	for i := 0; i < poolSize; i++ {
		go func(jobs <-chan *Job) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					res, err := poolWorkerFunc(job, writeInput.Append(), writeInput.Originals())
					logWork(res, job.filename, err)

					wg.Done()
				}
			}
		}(jobs)
	}
}

func poolWorker(job *Job, append, originals bool) ([]byte, error) {
	if append {
		//read from file and append to job.payload
	}

	filename, found := scan.Candidates(job.filename, files, config.Get().Extensions())
	if !found {
		return []byte{}, noFileErr
	}

	if len(job.payload.Tags()) == 0 {
		return []byte{}, skipFileErr
	}

	result, err := writecmd.WriteFile(
		filename,
		job.payload,
		originals,
	)

	return result, err
}

func logWork(result []byte, filename string, err error) {
	if err == skipFileErr {
		log.Debug("Skip", fmt.Sprintf("no payload found for `%s`", filename))
	} else if err == noFileErr {
		log.Debug("Skip", fmt.Sprintf("no files candidate for `%s`", filename))
	} else if err != nil {
		log.Failure("", err.Error())
	} else if len(result) != 0 {
		log.Success("Success", string(result))
	}
}
