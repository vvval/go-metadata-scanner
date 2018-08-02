package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/config"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/cmd/write"
	"github.com/vvval/go-metadata-scanner/log"
	"github.com/vvval/go-metadata-scanner/scan"
	syslog "log"
	"os"
	"path/filepath"
	"sync"
)

type Job struct {
	Directory     string
	Filename      string
	Payload       metadata.Payload
	SaveOriginals bool
	Append        bool
}

var (
	wg   sync.WaitGroup
	jobs chan *Job
)

const handlers = 20

var files []string

func init() {
	var writeCmd = &cobra.Command{
		Use:   "write",
		Short: "Read metadata from file and write to images",
		Long: `Read metadata from file and write to images.
Input file should be a CSV file with comma-separated fields.
First column should be reserved for file names, its name is omitted.
Other columns should be named as keywords in a config.yaml tagmap section provided
for proper mapping CSV data into appropriate metadata fields`,
		Run: execute,
	}

	files, err := scan.Dir(write.Input().Directory(), config.AppConfig().Extensions)
	if err != nil {
		syslog.Fatalln(err)
	}
	fmt.Printf("%+v\n", files)

	rootCmd.AddCommand(writeCmd)
	write.InitFlags(writeCmd)

	jobs = make(chan *Job)

	// worker pool
	for i := 0; i < handlers; i++ {
		go func(jobs chan *Job) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					res, err := work(job, &files)
					if err != nil {
						log.Failure("", err.Error())
					} else {
						log.Success("Processing", res)
					}

					wg.Done()
				}
			}
		}(jobs)
	}
}

func work(j *Job, files *[]string) (res string, err error) {
	if j.Append {
		//read from file and append to j.Payload
	}

	scan.Candidates(j.Filename, *files)

	result, err := write.WriteFile(
		filenameCandidates(j.Directory, j.Filename),
		j.Payload,
		j.SaveOriginals,
	)

	return string(result), err
}

func execute(cmd *cobra.Command, args []string) {
	input := write.Input()
	file, err := openFile(input.Filename())
	if err != nil {
		syslog.Fatalln(err)
	}
	defer file.Close()

	write.Read(file, input.Separator(), func(filename string, payload metadata.Payload) {
		wg.Add(1)
		jobs <- &Job{
			Directory:     input.Directory(),
			Filename:      filename,
			Payload:       payload,
			SaveOriginals: input.Originals(),
			Append:        input.Append(),
		}
	})

	wg.Wait()
	close(jobs)

	log.Log("Operation", "done")
}

func openFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func filenameCandidates(dir, name string) []string {
	return []string{filepath.Join(dir, name)}
}

//files, err := util.Read(input.Directory(), config.AppConfig().Extensions)
//if err != nil {
//	log.Fatalln(err)
//}
