package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd"
	"github.com/vvval/go-metadata-scanner/cmd/writeCommand"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/log"
	"github.com/vvval/go-metadata-scanner/metadata"
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
	wg          sync.WaitGroup
	jobs        chan *Job
	files       []string
	skipFileErr = errors.New("skipFileErr")
	noFileErr   = errors.New("noFileErr")
	writeInput  writeCommand.Input
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
	writeCommand.FillInput(cmd, &writeInput)

	jobs = make(chan *Job)
	initPool(handlers, jobs)
}

func appendScanPool(wg *sync.WaitGroup, poolSize int, files <-chan []string, callback func(files []string) ([]byte, error), jsonResult *[]map[string]interface{}) {
	for i := 0; i < poolSize; i++ {
		go func(files <-chan []string) {
			for {
				select {
				case chunk, ok := <-files:
					if !ok {
						return
					}

					//res, err := callback(chunk)

					//if err == nil {
					//	unmarshal := make([]map[string]interface{}, len(chunk))
					//	if err = json.Unmarshal(res, &unmarshal); err == nil {
					//		fmt.Printf("PING2 %+v\n", unmarshal)
					//		for _, elem := range unmarshal {
					//			*jsonResult = append(*jsonResult, elem)
					//		}
					//	} else {
					//		fmt.Printf("PING2 %+v\n", err.Error())
					//	}
					//} else {
					//
					//	fmt.Printf("PING3 %+v\n", err.Error())
					//}
					if res, err := callback(chunk); err == nil {
						unmarshal := make([]map[string]interface{}, len(chunk))
						if err = json.Unmarshal(res, &unmarshal); err == nil {
							for _, elem := range unmarshal {
								*jsonResult = append(*jsonResult, elem)
							}
						}
					}

					wg.Done()
				}
			}
		}(files)
	}
}

func writeHandler(cmd *cobra.Command, args []string) {
	if writeInput.Append() {
		jsonResult := make([]map[string]interface{}, len(files))
		scanFiles := make(chan []string)
		var scanWG sync.WaitGroup
		appendScanPool(&scanWG, scancmd.PoolSize, scanFiles, exec, &jsonResult)
		ReadJSON(writeInput.Directory(), config.Get().Extensions(), scancmd.PoolSize, &scanWG, scanFiles)

		wg.Wait()
		close(scanFiles)

		fmt.Printf("Output size: %+v\n", jsonResult)
		//read from file and append to job.payload
	}

	file, err := util.OpenReadonlyFile(writeInput.Filename())
	if err != nil {
		syslog.Fatalln(err)
	}
	defer file.Close()

	files = scanDir(writeInput.Directory(), config.Get().Extensions())

	writeCommand.ReadFile(file, writeInput.Separator(), func(filename string, payload metadata.Payload) {
		wg.Add(1)
		jobs <- &Job{filename, payload}
	})

	wg.Wait()
	close(jobs)

	log.Log("Writing", "done")
}

func scanDir(directory string, extensions []string) []string {
	result, err := scan.Dir(directory, extensions)
	if err != nil {
		syslog.Fatalln(err)
	}

	return result
}

func initPool(poolSize int, jobs <-chan *Job) {
	// worker pool
	for i := 0; i < poolSize; i++ {
		go func(jobs <-chan *Job) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					res, err := work(job, writeInput)
					logWork(res, job.filename, err)

					wg.Done()
				}
			}
		}(jobs)
	}
}

func work(job *Job, input writeCommand.Input) (res []byte, err error) {
	if input.Append() {
		//read from file and append to job.payload
	}

	filename, found := scan.Candidates(job.filename, files, config.Get().Extensions())
	if !found {
		return []byte{}, noFileErr
	}

	if len(job.payload.Tags()) == 0 {
		return []byte{}, skipFileErr
	}

	result, err := writeCommand.WriteFile(
		filename,
		job.payload,
		input.Originals(),
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
