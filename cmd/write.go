package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/cmd/writer"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/log"
	"github.com/vvval/go-metadata-scanner/scan"
	syslog "log"
	"os"
	"sync"
)

type Job struct {
	filename string
	payload  metadata.Payload
}

var (
	wg   sync.WaitGroup
	jobs chan *Job
)

const handlers int = 20

var files []string

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
	writer.InitFlags(cmd)

	files = scanDir()
	initPool(handlers)
}

func scanDir() []string {
	result, err := scan.Dir(writer.Input().Directory(), config.Get().Extensions())
	if err != nil {
		syslog.Fatalln(err)
	}

	return result
}

func initPool(poolSize int) {
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

					res, err := work(job)
					if err != nil {
						log.Failure("", err.Error())
					} else {
						log.Success("Writing", res)
					}

					wg.Done()
				}
			}
		}(jobs)
	}
}

func work(job *Job) (res string, err error) {
	input := writer.Input()

	if input.Append() {
		//read from file and append to job.payload
	}

	filename, found := scan.Candidates(job.filename, files, config.Get().Extensions())

	if !found {
		return "", fmt.Errorf("no file candidate for `%s`", job.filename)
	}

	result, err := writer.WriteFile(
		filename,
		job.payload,
		input.Originals(),
	)

	return string(result), err
}

func writeHandler(cmd *cobra.Command, args []string) {
	input := writer.Input()
	file, err := openFile(input.Filename())
	if err != nil {
		syslog.Fatalln(err)
	}
	defer file.Close()

	writer.Read(file, input.Separator(), func(filename string, payload metadata.Payload) {
		wg.Add(1)
		jobs <- &Job{filename, payload}
	})

	wg.Wait()
	close(jobs)

	log.Log("Writing", "done")
}

func openFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}
