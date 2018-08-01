package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/bwrite"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/wolfy-j/goffli/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

//type Job struct {
//	Line          metadata.Line
//	Directory     string
//	Filename      string
//	SaveOriginals bool
//	Append        bool
//}

var (
	wg     sync.WaitGroup
	jobs   chan bwrite.Job
	status chan string
	errors chan error
)

const handlers = 20

func init() {
	var bwriteCmd = &cobra.Command{
		Use:   "bwrite",
		Short: "Read metadata from file and write to images",
		Long: `Read metadata from file and write to images. Input file should be a CSV file with comma-separated fields.
First column should be reserved for file names, its name is omitted.
Other columns should be named as keywords in a config.yaml tagmap section provided
for proper mapping CSV data into appropriate metadata fields`,
		Run: bulkwrite,
	}

	rootCmd.AddCommand(bwriteCmd)
	bwrite.InitFlags(bwriteCmd)

	jobs = make(chan bwrite.Job)
	status = make(chan string)
	errors = make(chan error)

	// worker pool
	for i := 0; i < handlers; i++ {
		go func(jobs <-chan bwrite.Job, status chan string, errors chan error) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					res, err := work(job)
					if err != nil {
						errors <- err
					} else {
						status <- res
					}

					wg.Done()
				}
			}
		}(jobs, status, errors)
	}
}

func work(j bwrite.Job) (res string, err error) {
	if j.Append {
		//read from file and append to j.Line
	}
	result, err := bwrite.WriteFile(
		filenameCandidates(j.Directory, j.Filename),
		j.Line,
		j.SaveOriginals,
	)

	return string(result), err
}

func bulkwrite(cmd *cobra.Command, args []string) {
	go func() {
		for {
			select {
			case _, ok := <-status:
				if !ok {
					return
				}
			case err, ok := <-errors:
				if ok {
					util.LogError("", err.Error())
				}
			}
		}
	}()

	input := bwrite.Input()
	file, err := openFile(input.Filename())
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	reader, err := bwrite.Reader(file, input.Separator())
	if err != nil {
		log.Fatalln(err)
	}

	//bwrite.ScanFile(reader, wg, jobs, input)

	var columnsLineFound bool
	var columns map[int]string

	var i int
	for {
		i++
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			utils.Printf("<red>Error on reading line %d: %s</reset>\n", i, err)
			continue
		}

		if !columnsLineFound {
			columnsLineFound = true
			columns = bwrite.MapColumns(line)

			var cols []string
			for _, col := range columns {
				cols = append(cols, col)
			}

			utils.Log("Mapped columns are:", cols...)
			continue
		}

		if skipLine(line) {
			continue
		}

		wg.Add(1)
		job := bwrite.Job{
			Line:          bwrite.MapLineToColumns(columns, line),
			Directory:     input.Directory(),
			Filename:      line[0],
			SaveOriginals: input.Originals(),
			Append:        input.Append(),
		}
		jobs <- job
	}

	//files, err := util.Scan(input.Directory(), config.AppConfig().Extensions)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	wg.Wait()
	close(jobs)
	close(status)
	close(errors)

	utils.Log("Operation", "done")
}

func openFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func skipLine(line []string) bool {
	return len(line) == 0 || len(line[0]) == 0
}

func filenameCandidates(dir, name string) []string {
	return []string{filepath.Join(dir, name)}
}
