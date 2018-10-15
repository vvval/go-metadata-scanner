package writecmd

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/util/log"
	"sync"
)

func CreatePool(
	wg *sync.WaitGroup,
	poolSize int,
	jobs <-chan *Job,
	poolWorkerFunc func(job *Job) ([]byte, error),
) {
	// worker pool
	for i := 0; i < poolSize; i++ {
		go func(jobs <-chan *Job) {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					_, err := poolWorkerFunc(job)
					if err != nil {
						logError(job.Filename(), err)
					}
					wg.Done()
				}
			}
		}(jobs)
	}
}

func logError(filename string, err error) {
	if err == skipFileErr {
		log.Debug("Skip file", fmt.Sprintf("no payload found for `%s`", filename))
	} else if err == noFileErr {
		log.Debug("Skip file", fmt.Sprintf("no files candidate for `%s`", filename))
	} else {
		log.Failure("Write error", err.Error())
	}
}
