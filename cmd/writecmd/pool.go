package writecmd

import (
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
						log.Failure("Write error", err.Error())
					}
					wg.Done()
				}
			}
		}(jobs)
	}
}
