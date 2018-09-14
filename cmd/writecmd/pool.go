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
	poolWorkerFunc func(job *Job, append, originals bool) ([]byte, error),
	append, originals bool,
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

					_, err := poolWorkerFunc(job, append, originals)
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
	if err == SkipFileErr {
		log.Debug("Skip", fmt.Sprintf("no payload found for `%s`", filename))
	} else if err == NoFileErr {
		log.Debug("Skip", fmt.Sprintf("no files candidate for `%s`", filename))
	} else {
		log.Failure("Write error", err.Error())
	}
}
