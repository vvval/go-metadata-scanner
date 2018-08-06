package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/cmd/writeCommand"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/log"
	"github.com/vvval/go-metadata-scanner/scan"
	"github.com/vvval/go-metadata-scanner/util"
	syslog "log"
	"sync"
)

type Job struct {
	filename string
	payload  metadata.Payload
}

var (
	wg       sync.WaitGroup
	jobs     chan *Job
	files    []string
	skipFile = errors.New("skipFile")
	noFile   = errors.New("noFile")
	input    = writeCommand.Input{}
)

const handlers int = 20

func init() {
	var cmd = &cobra.Command{
		Use:   "write",
		Short: "Read metadata from file and writer to images",
		Long: `Read metadata from file and writer to images.
Input file should be a CSV file with comma-separated fields (or pass custom separator via "s" flag).
First column should be reserved for file names, its name is omitted.
Other columns should be named as keywords in a dict.yaml maps section provided
for proper mapping CSV data into appropriate metadata fields`,
		Run: writeHandler,
	}

	rootCmd.AddCommand(cmd)
	writeCommand.FillInput(cmd, &input)

	initPool(cmd, handlers)
}

func writeHandler(cmd *cobra.Command, args []string) {
	file, err := util.OpenReadonlyFile(input.Filename())
	if err != nil {
		syslog.Fatalln(err)
	}
	defer file.Close()

	files = scanDir(input)

	writeCommand.ReadFile(file, input.Separator(), func(filename string, payload metadata.Payload) {
		wg.Add(1)
		jobs <- &Job{filename, payload}
	})

	wg.Wait()
	close(jobs)

	log.Log("Writing", "done")
}

func scanDir(input writeCommand.Input) []string {
	result, err := scan.Dir(input.Directory(), config.Get().Extensions())
	if err != nil {
		syslog.Fatalln(err)
	}

	return result
}

func initPool(cmd *cobra.Command, poolSize int) {
	jobs = make(chan *Job)

	// worker pool
	for i := 0; i < poolSize; i++ {
		go func(jobs chan *Job) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					res, err := work(job, input)
					logWork(res, job.filename, err)

					wg.Done()
				}
			}
		}(jobs)
	}
}

func work(job *Job, input writeCommand.Input) (res string, err error) {
	if input.Append() {
		//read from file and append to job.payload
	}

	filename, found := scan.Candidates(job.filename, files, config.Get().Extensions())
	if !found {
		return "", noFile
	}

	if len(job.payload.Tags()) == 0 {
		return "", skipFile
	}

	result, err := writeCommand.WriteFile(
		filename,
		job.payload,
		input.Originals(),
	)

	return string(result), err
}

func logWork(res, filename string, err error) {
	if err == skipFile {
		log.Debug("Skip", fmt.Sprintf("no payload found for `%s`", filename))
	} else if err == noFile {
		log.Debug("Skip", fmt.Sprintf("no files candidate for `%s`", filename))
	} else if err != nil {
		log.Failure("", err.Error())
	} else if len(res) != 0 {
		log.Success("Success", res)
	}
}
