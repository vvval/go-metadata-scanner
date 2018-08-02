package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/bwrite"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/wolfy-j/goffli/utils"
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
	wg   sync.WaitGroup
	jobs chan *bwrite.Job
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

	jobs = make(chan *bwrite.Job)

	// worker pool
	for i := 0; i < handlers; i++ {
		go func(jobs chan *bwrite.Job) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					res, err := work(job)
					if err != nil {
						util.Error("", err.Error())
					} else {
						utils.Printf("<green>%s</reset>", res)
					}

					wg.Done()
				}
			}
		}(jobs)
	}
}

func work(j *bwrite.Job) (res string, err error) {
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

	bwrite.ScanFile(reader, func(filename string, line metadata.Line) {
		wg.Add(1)
		jobs <- &bwrite.Job{
			Directory:     input.Directory(),
			Filename:      filename,
			Line:          line,
			SaveOriginals: input.Originals(),
			Append:        input.Append(),
		}
	})

	wg.Wait()
	close(jobs)

	utils.Log("Operation", "done")
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

//files, err := util.Scan(input.Directory(), config.AppConfig().Extensions)
//if err != nil {
//	log.Fatalln(err)
//}
